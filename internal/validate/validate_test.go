package validate

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

type sampleServer struct {
	Address string   `mapstructure:"address" validate:"required,ipv4"`
	Port    int      `mapstructure:"port" validate:"required,gte=1024,lte=65535"`
	Origins []string `mapstructure:"allow_origins" validate:"omitempty,dive,oneof=* https://example.com"`
}

type sampleConfig struct {
	LogLevel string       `mapstructure:"log_level" validate:"omitempty,oneof=debug info warn error"`
	Server   sampleServer `mapstructure:"server"`
	Email    string       `mapstructure:"email" validate:"omitempty,email"`
}

type sampleWithoutTags struct {
	ConfigFile string `validate:"required"`
}

// The formatted error must be readable on its own: one line per violation, each
// naming the configuration key instead of the Go struct path.
func TestStruct_FormatsOneLinePerViolation(t *testing.T) {
	err := Struct(sampleConfig{
		LogLevel: "verbose",
		Server:   sampleServer{Address: "not-an-ip", Port: 70000},
		Email:    "not-an-email",
	})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	lines := strings.Split(err.Error(), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected one line per violation (4), got %d:\n%s", len(lines), err)
	}

	for _, expected := range []string{
		"- log_level must be one of [debug info warn error]",
		"- server.address must be a valid IPv4 address",
		"- server.port must be 65,535 or less",
		"- email must be a valid email address",
	} {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected %q in error:\n%s", expected, err)
		}
	}
	if strings.Contains(err.Error(), "sampleConfig") {
		t.Errorf("expected configuration keys instead of Go type names:\n%s", err)
	}
}

func TestStruct_ReportsMissingRequiredField(t *testing.T) {
	err := Struct(sampleWithoutTags{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if !strings.Contains(err.Error(), "- configFile is a required field") {
		t.Fatalf("expected camelCase fallback field name, got: %s", err)
	}
}

func TestStruct_ReturnsNilForValidConfig(t *testing.T) {
	err := Struct(sampleConfig{
		LogLevel: "info",
		Server:   sampleServer{Address: "127.0.0.1", Port: 8156, Origins: []string{"*"}},
		Email:    "user@example.com",
	})
	if err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestPretty_PassesThroughNonValidationErrors(t *testing.T) {
	original := errors.New("boom")

	if got := pretty(original); got != original {
		t.Fatalf("expected original error, got %v", got)
	}
}

func TestVar_FormatsSingleValueViolation(t *testing.T) {
	if err := Var("not-a-cidr", "cidr"); err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if err := Var("10.0.0.0/8", "cidr"); err != nil {
		t.Fatalf("expected valid CIDR, got: %v", err)
	}
}

func TestRegisterValidation_UsesSharedInstance(t *testing.T) {
	if err := RegisterValidation("always_fails", func(validator.FieldLevel) bool { return false }); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	type target struct {
		Name string `mapstructure:"name" validate:"always_fails"`
	}

	got := Struct(target{Name: "x"})
	if got == nil || !strings.Contains(got.Error(), "- name failed the \"always_fails\" rule") {
		t.Fatalf("expected formatted custom rule error, got: %v", got)
	}
}

// Custom rules must render the same readable message as built-in ones.
func TestRegisterRule_RendersTranslatedMessage(t *testing.T) {
	if err := RegisterRule("must_be_lowercase", "{0} must be lowercase", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == strings.ToLower(fl.Field().String())
	}); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	type target struct {
		Name string `mapstructure:"name" validate:"must_be_lowercase"`
	}

	got := Struct(target{Name: "UPPER"})
	if got == nil {
		t.Fatal("expected validation error, got nil")
	}
	if want := "- name must be lowercase"; got.Error() != want {
		t.Fatalf("got %q, want %q", got.Error(), want)
	}

	if err := Struct(target{Name: "lower"}); err != nil {
		t.Fatalf("expected valid value, got: %v", err)
	}
}

func TestIsConfigFile(t *testing.T) {
	cases := []struct {
		name string
		path string
		ok   bool
	}{
		{name: "yaml", path: "./config/config.yaml", ok: true},
		{name: "yml", path: "./config/config.yml", ok: true},
		{name: "uppercase extension", path: "./config/config.YAML", ok: true},
		{name: "txt", path: "./config/config.txt", ok: false},
		{name: "dot", path: ".", ok: false},
		{name: "root", path: string(filepath.Separator), ok: false},
		{name: "traversal", path: "../config/config.yaml", ok: false},
		{name: "nested traversal", path: "./config/../../secret.yaml", ok: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Var(tc.path, configFileTag) == nil; got != tc.ok {
				t.Fatalf("isConfigFile(%q) valid = %v, want %v", tc.path, got, tc.ok)
			}
		})
	}
}

func TestIsValidEnvKey(t *testing.T) {
	cases := []struct {
		key string
		ok  bool
	}{
		{key: "FOO", ok: true},
		{key: "_PRIVATE", ok: true},
		{key: "A1_B2", ok: true},
		{key: "1INVALID", ok: false},
		{key: "FOO=BAR", ok: false},
		{key: "WITH SPACE", ok: false},
		{key: "", ok: false},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			if got := Var(tc.key, envKeyTag) == nil; got != tc.ok {
				t.Fatalf("isValidEnvKey(%q) valid = %v, want %v", tc.key, got, tc.ok)
			}
		})
	}
}

func TestAllowAllOriginsOrOrigin(t *testing.T) {
	cases := []struct {
		origin string
		ok     bool
	}{
		{origin: "*", ok: true},
		{origin: "https://example.com", ok: true},
		{origin: "http://localhost:5173", ok: true},
		{origin: "https://example.com/", ok: false},
		{origin: "https://example.com/path", ok: false},
		{origin: "https://example.com?q=1", ok: false},
		{origin: "example.com", ok: false},
		{origin: "not a url", ok: false},
	}

	for _, tc := range cases {
		t.Run(tc.origin, func(t *testing.T) {
			if got := Var(tc.origin, corsOriginTag) == nil; got != tc.ok {
				t.Fatalf("allowAllOriginsOrOrigin(%q) valid = %v, want %v", tc.origin, got, tc.ok)
			}
		})
	}
}

func TestConfigFieldName(t *testing.T) {
	cases := []struct {
		name  string
		field sampleServer
		index int
		want  string
	}{
		{name: "mapstructure tag", field: sampleServer{}, index: 0, want: "address"},
		{name: "fallback to field name", field: sampleServer{}, index: 1, want: "port"},
	}

	for _, tc := range cases {
		got := configFieldName(reflect.TypeOf(tc.field).Field(tc.index))
		if got != tc.want {
			t.Errorf("%s: configFieldName() = %q, want %q", tc.name, got, tc.want)
		}
	}
}
