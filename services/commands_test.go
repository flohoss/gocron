package services

import (
	"strings"
	"testing"

	"github.com/flohoss/gocron/config"
)

func testTerminalSettings() config.TerminalSettings {
	settings := config.TerminalSettings{
		AllowAllCommands: false,
		AllowedCommands: map[string]config.AllowedCommands{
			"docker": {
				Args: []string{"ps", "version"},
			},
			"cat": {
				Args: []string{"/config/config.yaml"},
			},
			"echo": {},
		},
	}
	settings.Hydrate()
	return settings
}

func TestExecute_EmptyCommandRejected(t *testing.T) {
	_, err := execute("   ", testTerminalSettings())
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
	if err.Error() != "empty command" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_DisallowedCommandRejected(t *testing.T) {
	_, err := execute("whoami", testTerminalSettings())
	if err == nil {
		t.Fatal("expected error for disallowed command, got nil")
	}
	if !strings.Contains(err.Error(), `command "whoami" is not allowed`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_DisallowedArgumentRejected(t *testing.T) {
	_, err := execute("docker images", testTerminalSettings())
	if err == nil {
		t.Fatal("expected error for disallowed argument, got nil")
	}
	if !strings.Contains(err.Error(), `argument "images" is not allowed for command "docker"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_CommandWithNoArgsRejectedWhenArgsProvided(t *testing.T) {
	_, err := execute("echo hello", testTerminalSettings())
	if err == nil {
		t.Fatal("expected error for command arguments, got nil")
	}
	if !strings.Contains(err.Error(), `command "echo" does not allow any arguments`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_AllowAllArgsCommandAccepted(t *testing.T) {
	settings := config.TerminalSettings{
		AllowAllCommands: false,
		AllowedCommands: map[string]config.AllowedCommands{
			"echo": {AllowAllArgs: true},
		},
	}

	out, err := execute("echo hello", settings)
	if err != nil {
		t.Fatalf("expected command to succeed, got error: %v", err)
	}
	if !strings.Contains(out, "hello") {
		t.Fatalf("unexpected output: %q", out)
	}
}
