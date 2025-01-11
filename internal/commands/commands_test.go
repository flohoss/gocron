package commands

import (
	"os"
	"testing"
)

func TestExtractVariable(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("lala", "resolved")
	os.Setenv("HOME", "/home/user")
	os.Setenv("VAR1", "path")
	os.Setenv("VAR2", "to/resource")

	tests := []struct {
		input    string
		expected string
	}{
		// Basic test case: variable in the middle
		{"/lili/${lala}/Lulu", "/lili/resolved/Lulu"},
		// No variables
		{"plain-text", "plain-text"},
		// Unresolved variable
		{"/lili/${undefined}/Lulu", "/lili/${undefined}/Lulu"},
		// Variable at the start
		{"${HOME}/docs", "/home/user/docs"},
		// Multiple variables in the same string
		{"/${VAR1}/${VAR2}/file", "/path/to/resource/file"},
		// Nested variables (should be treated as unresolved)
		{"/${${VAR1}}/suffix", "/${${VAR1}}/suffix"},
		// Missing closing brace
		{"/${lala/missing", "/${lala/missing"},
		// Empty variable name
		{"/${}/suffix", "/${}/suffix"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ExtractVariable(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractVariable(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPrepareResticCommand(t *testing.T) {
	// Set environment variables
	os.Setenv("RESTIC_POLICY", "--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 75")
	os.Setenv("BASE_REPOSITORY", "rclone:pcloud:Server/Backups")

	// Define the command
	command := "restic -r ${BASE_REPOSITORY}/directus forget ${RESTIC_POLICY} --prune"

	// Prepare the command
	program, args := PrepareCommand(command)

	// Expected values
	expectedProgram := "restic"
	expectedArgs := []string{
		"-r", "rclone:pcloud:Server/Backups/directus",
		"forget",
		"--keep-daily", "7",
		"--keep-weekly", "5",
		"--keep-monthly", "12",
		"--keep-yearly", "75",
		"--prune",
	}

	// Check program
	if program != expectedProgram {
		t.Errorf("Expected program '%s', but got '%s'", expectedProgram, program)
	}

	// Check arguments
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d arguments, but got %d", len(expectedArgs), len(args))
	}
	for i, arg := range args {
		if arg != expectedArgs[i] {
			t.Errorf("Argument %d: expected '%s', but got '%s'", i, expectedArgs[i], arg)
		}
	}
}
