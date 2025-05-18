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

func TestPrepareResticCommand(t *testing.T) {
	// Set environment variables
	os.Setenv("RESTIC_POLICY", "--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 75")
	os.Setenv("BASE_REPOSITORY", "rclone:pcloud:Server/Backups")
	os.Setenv("BACKUP_PATH", "/app/config/test")

	// Define the commands to test
	testCases := []struct {
		command         string
		expectedProgram string
		expectedArgs    []string
	}{
		{
			command:         "restic -r ${BASE_REPOSITORY}/directus forget ${RESTIC_POLICY} --prune",
			expectedProgram: "restic",
			expectedArgs: []string{
				"-r", "rclone:pcloud:Server/Backups/directus",
				"forget",
				"--keep-daily", "7",
				"--keep-weekly", "5",
				"--keep-monthly", "12",
				"--keep-yearly", "75",
				"--prune",
			},
		},
		{
			command:         "find '/mnt/cache/docker/appdata/plex/Library/Application Support/Plex Media Server/Cache/PhotoTranscoder' -name \"*.jpg\" -type f -mtime +5 -print -delete",
			expectedProgram: "find",
			expectedArgs: []string{
				"/mnt/cache/docker/appdata/plex/Library/Application Support/Plex Media Server/Cache/PhotoTranscoder",
				"-name", "*.jpg",
				"-type", "f",
				"-mtime", "+5",
				"-print",
				"-delete",
			},
		},
		{
			command:         "echo 'Hello, World!'",
			expectedProgram: "echo",
			expectedArgs:    []string{"Hello, World!"},
		},
		{
			command:         "ls -la /tmp",
			expectedProgram: "ls",
			expectedArgs:    []string{"-la", "/tmp"},
		},
		{
			command:         "grep 'error' /var/log/syslog",
			expectedProgram: "grep",
			expectedArgs:    []string{"error", "/var/log/syslog"},
		},
		{
			command:         "tar -czvf archive.tar.gz /home/user",
			expectedProgram: "tar",
			expectedArgs:    []string{"-czvf", "archive.tar.gz", "/home/user"},
		},
		{
			command:         "ping -c 4 google.com",
			expectedProgram: "ping",
			expectedArgs:    []string{"-c", "4", "google.com"},
		},
		{
			command:         "rm -r -f ${BACKUP_PATH}/*",
			expectedProgram: "rm",
			expectedArgs:    []string{"-r", "-f", "/app/config/test/*"},
		},
	}

	for _, testCase := range testCases {
		fmt.Printf("Testing command: %s\n", testCase.command)
		program, args := PrepareCommand(testCase.command)
		fmt.Printf("Expected Program: %s, Got: %s\n", testCase.expectedProgram, program)
		fmt.Printf("Expected Args: %v, Got: %v\n", testCase.expectedArgs, args)

		if program != testCase.expectedProgram {
			t.Errorf("Expected program '%s', but got '%s'", testCase.expectedProgram, program)
		}

		if len(args) != len(testCase.expectedArgs) {
			t.Fatalf("Expected %d arguments, but got %d", len(testCase.expectedArgs), len(args))
		}

		for i, arg := range args {
			if arg != testCase.expectedArgs[i] {
				t.Errorf("Argument %d: expected '%s', but got '%s'", i, testCase.expectedArgs[i], arg)
			}
		}
	}
}
