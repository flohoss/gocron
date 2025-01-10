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
