package commands

import (
	"fmt"
	"os"
	"testing"
)

func TestExtractVariable(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("lala", "resolved")
	os.Setenv("HOME", "/home/user")
	os.Setenv("VAR1", "path")
	os.Setenv("VAR2", "to/resource")
	os.Setenv("EMPTY", "")
	os.Setenv("SPECIAL_CHARS", "value!@#")
	os.Setenv("NUMERIC", "12345")

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
		// Empty variable value
		{"/prefix/${EMPTY}/suffix", "/prefix//suffix"},
		// Variable with special characters
		{"/prefix/${SPECIAL_CHARS}/suffix", "/prefix/value!@#/suffix"},
		// Variable with numeric values
		{"/data/${NUMERIC}/file", "/data/12345/file"},
		// Variable at the end
		{"/home/user/${VAR1}", "/home/user/path"},
		// Multiple instances of the same variable
		{"/${VAR1}/${VAR1}/again", "/path/path/again"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			fmt.Printf("Testing input: %s\n", tt.input)
			result := ExtractVariable(tt.input)
			fmt.Printf("Expected: %s, Got: %s\n", tt.expected, result)
			if result != tt.expected {
				t.Errorf("ExtractVariable(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
