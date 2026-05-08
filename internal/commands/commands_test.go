package commands

import (
	"strings"
	"testing"
)

func TestExecuteCommand_ReturnsOutputOnSuccess(t *testing.T) {
	out, err := ExecuteCommand("printf hello")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "hello" {
		t.Fatalf("unexpected output: got %q want %q", out, "hello")
	}
}

func TestExecuteCommand_ReturnsNoOutputWhenCommandSilent(t *testing.T) {
	out, err := ExecuteCommand("true")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "No output" {
		t.Fatalf("unexpected output: got %q want %q", out, "No output")
	}
}

func TestExecuteCommand_ReturnsErrorWhenCommandFails(t *testing.T) {
	out, err := ExecuteCommand("sh -c 'echo fail >&2; exit 1'")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(out, "fail") {
		t.Fatalf("unexpected output: %q", out)
	}
}
