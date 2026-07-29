package commands

import (
	"context"
	"strings"
	"testing"
	"time"
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

func TestExecuteCommandWithContext_RespectsTimeout(t *testing.T) {
	out, err := ExecuteCommandWithContext(context.Background(), "sleep 5; echo done", 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error, got: %v", err)
	}
	_ = out
}

func TestExecuteCommandWithContext_CompletesWithinTimeout(t *testing.T) {
	out, err := ExecuteCommandWithContext(context.Background(), "printf fast", 1*time.Second)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "fast" {
		t.Fatalf("unexpected output: got %q want %q", out, "fast")
	}
}

func TestExecuteCommandWithContext_NoTimeoutBehavesLikeExecuteCommand(t *testing.T) {
	out, err := ExecuteCommandWithContext(context.Background(), "sh -c 'echo fail >&2; exit 1'", 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(out, "fail") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRetryLoop_RetriesFailingCommand(t *testing.T) {
	cmd := "exit 1"
	var attempts int
	var lastErr error
	retries := 2

	for attempt := 0; attempt <= retries; attempt++ {
		attempts++
		_, err := ExecuteCommandWithContext(context.Background(), cmd, 0)
		if err == nil {
			break
		}
		lastErr = err
	}

	if attempts != retries+1 {
		t.Fatalf("expected %d attempts, got %d", retries+1, attempts)
	}
	if lastErr == nil {
		t.Fatal("expected final error after all retries exhausted")
	}
}

func TestRetryLoop_StopsOnFirstSuccess(t *testing.T) {
	cmd := "true"
	var attempts int
	retries := 3

	for attempt := 0; attempt <= retries; attempt++ {
		attempts++
		_, err := ExecuteCommandWithContext(context.Background(), cmd, 0)
		if err == nil {
			break
		}
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt for successful command, got %d", attempts)
	}
}

func TestRetryLoop_RetriesUntilSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	marker := tmpDir + "/attempts"

	cmd := "sh -c 'f=" + marker + "; n=$(cat \"$f\" 2>/dev/null || echo 0); n=$((n+1)); echo $n > \"$f\"; if [ $n -lt 3 ]; then exit 1; fi; echo success'"

	var attempts int
	var lastOut string
	var lastErr error
	retries := 5

	for attempt := 0; attempt <= retries; attempt++ {
		attempts++
		out, err := ExecuteCommandWithContext(context.Background(), cmd, 0)
		lastOut = out
		lastErr = err
		if err == nil {
			break
		}
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts (2 retries) to reach success, got %d", attempts)
	}
	if lastErr != nil {
		t.Fatalf("expected success after retries, got error: %v", lastErr)
	}
	if !strings.Contains(lastOut, "success") {
		t.Fatalf("expected success output, got %q", lastOut)
	}
}
