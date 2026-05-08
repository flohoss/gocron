package cli

import (
	"os"
	"strings"
	"testing"
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
	dir := t.TempDir()

	opts, err := Parse([]string{"-config", dir})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.HasSuffix(opts.ConfigFolder, string(os.PathSeparator)) {
		t.Fatalf("expected config path to end with path separator, got %q", opts.ConfigFolder)
	}
}

func TestParse_WithInvalidFlag(t *testing.T) {
	_, err := Parse([]string{"-unknown-flag"})
	if err == nil {
		t.Fatal("expected parse error for unknown flag, got nil")
	}
}

func TestParse_EmptyConfig_UsesDefault(t *testing.T) {
	// When -config flag is not provided, default config folder is used
	opts, err := Parse([]string{})
	if err != nil {
		t.Fatalf("expected no error with default config, got: %v", err)
	}
	if opts.ConfigFolder == "" {
		t.Fatal("expected non-empty default config folder")
	}
}

func TestNormalizeDirPath_DotBecomesSlash(t *testing.T) {
	got := normalizeDirPath(".")
	if got != "."+string(os.PathSeparator) {
		t.Fatalf("expected ./ for dot input, got %q", got)
	}
}
