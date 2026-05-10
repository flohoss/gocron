package cli

import (
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestParse_WithVersionSkipsValidation(t *testing.T) {
	opts, err := Parse([]string{"-version", "-config", "/path/that/does/not/exist"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !opts.ShowVersion {
		t.Fatal("expected ShowVersion to be true")
	}
}

func TestParse_WithValidConfigPath(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "custom-config.yaml")

	opts, err := Parse([]string{"-config", configFile})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if opts.ConfigFile != configFile {
		t.Fatalf("expected config file path %q, got %q", configFile, opts.ConfigFile)
	}
}

func TestParse_WithInvalidFlag(t *testing.T) {
	_, err := Parse([]string{"-unknown-flag"})
	if err == nil {
		t.Fatal("expected parse error for unknown flag, got nil")
	}
}

func TestParse_WithInvalidConfigExtension(t *testing.T) {
	_, err := Parse([]string{"-config", "./config/config.txt"})
	if err == nil {
		t.Fatal("expected parse error for invalid config extension, got nil")
	}
}

func TestParse_WithPathTraversalConfig(t *testing.T) {
	_, err := Parse([]string{"-config", "../config/config.yaml"})
	if err == nil {
		t.Fatal("expected parse error for path traversal config, got nil")
	}
}

func TestParse_EmptyConfig_UsesDefault(t *testing.T) {
	// When -config flag is not provided, default config file is used
	opts, err := Parse([]string{})
	if err != nil {
		t.Fatalf("expected no error with default config, got: %v", err)
	}
	if opts.ConfigFile == "" {
		t.Fatal("expected non-empty default config file")
	}
}

func TestNormalizeFilePath_CleansPath(t *testing.T) {
	got := normalizeFilePath("./config/../config/config.yaml")
	want := filepath.Clean("./config/../config/config.yaml")
	if got != want {
		t.Fatalf("expected cleaned path %q, got %q", want, got)
	}
}

func TestValidateConfigFile(t *testing.T) {
	cases := []struct {
		name string
		path string
		ok   bool
	}{
		{name: "yaml", path: "./config/config.yaml", ok: true},
		{name: "yml", path: "./config/config.yml", ok: true},
		{name: "dot", path: ".", ok: false},
		{name: "txt", path: "./config/config.txt", ok: false},
	}

	v := validator.New()
	if err := v.RegisterValidation("config_file", validateConfigFile); err != nil {
		t.Fatalf("unexpected register validation error: %v", err)
	}

	for _, tc := range cases {
		err := v.Var(tc.path, "config_file")
		if tc.ok && err != nil {
			t.Fatalf("case %q: expected success, got error %v", tc.name, err)
		}
		if !tc.ok && err == nil {
			t.Fatalf("case %q: expected error, got success", tc.name)
		}
	}
}
