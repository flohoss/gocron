package services

import (
	"os"
	"strings"
	"testing"

	"gitlab.unjx.de/flohoss/gocron/config"
)

// Helper function to create a mock terminal configuration
// This now returns the full Terminal struct.
func createTerminalConfig(allowAll bool) config.TerminalSettings {
	return config.TerminalSettings{
		AllowAllCommands: allowAll,
		AllowedCommands: []config.AllowedCommands{
			{Command: "cat", Args: []string{"/config/config.yaml"}},
			{Command: "ls", Args: []string{"-la"}},
			{Command: "docker", Args: []string{"ps"}},
			{Command: "hostname", Args: []string{}},
		},
	}
}

func TestExecute_Success(t *testing.T) {
	terminalConfig := createTerminalConfig(false)

	// Prepare a temporary file for the 'cat' command test
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file afterward
	tempFile.WriteString("terminal:\n  allowed_commands:\n")
	tempFile.Close()

	// Temporarily update the config to use the new temp file for the test
	terminalConfig.AllowedCommands[0].Args[0] = tempFile.Name()

	testCases := []struct {
		name             string
		command          string
		expectedContains string
	}{
		{"ls command", "ls -la", "total"},
		{"docker command", "docker ps", "CONTAINER ID"},
		{"hostname command", "hostname", ""},
		// Now that we've prepared a file, this test should pass.
		{"cat command", "cat " + tempFile.Name(), "terminal:"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execute(tc.command, terminalConfig)
			if err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}
			if tc.expectedContains != "" && !strings.Contains(output, tc.expectedContains) {
				t.Errorf("expected output to contain %q, but got: %q", tc.expectedContains, output)
			}
		})
	}
}

func TestExecute_Failure(t *testing.T) {
	terminalConfig := createTerminalConfig(false)

	testCases := []struct {
		name          string
		command       string
		expectedError string
	}{
		// Command not in the whitelist
		{"disallowed command", "rm -rf /", `command "rm" is not allowed`},
		{"disallowed command with no args", "uptime", `command "uptime" is not allowed`},

		// Arguments not allowed for the command
		{"disallowed ls arg", "ls -l", `argument "-l" is not allowed for command "ls"`},
		{"disallowed docker arg", "docker info", `argument "info" is not allowed for command "docker"`},
		{"disallowed cat arg", "cat /etc/passwd", `argument "/etc/passwd" is not allowed for command "cat"`},

		// No arguments allowed for command, but some were provided
		{"unexpected args for hostname", "hostname -s", `command "hostname" does not allow any arguments`},

		// Empty command string
		{"empty command", "", "empty command"},
		{"whitespace command", "  ", "empty command"},

		// Non-existent file for 'cat' command
		{"cat non-existent file", "cat /config/config.yaml", "No such file or directory"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execute(tc.command, terminalConfig)
			if err == nil {
				t.Fatalf("expected an error, but got none. Output: %q", output)
			}

			// Check if either the error message or the output contains the expected string.
			// This handles cases like 'cat' where the error is written to stdout/stderr.
			if !strings.Contains(err.Error(), tc.expectedError) && !strings.Contains(output, tc.expectedError) {
				t.Errorf("expected error message or output to contain %q, but got error: %v and output: %q", tc.expectedError, err, output)
			}
		})
	}
}

func TestExecute_AllowAll(t *testing.T) {
	terminalConfig := createTerminalConfig(true) // Set allowAll to true

	testCases := []struct {
		name             string
		command          string
		expectedContains string
	}{
		{"allowed command with args", "ls -la", "total"},
		{"disallowed command", "uptime", ""},
		{"command with no output", "true", "No output"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execute(tc.command, terminalConfig)

			if err != nil {
				t.Fatalf("expected no error for command %q, but got: %v", tc.command, err)
			}

			// For successful commands, check the output
			if tc.expectedContains != "" && !strings.Contains(output, tc.expectedContains) {
				t.Errorf("expected output to contain %q, but got: %q", tc.expectedContains, output)
			}
		})
	}
}

func TestExecute_NoOutput(t *testing.T) {
	terminalConfig := createTerminalConfig(false)
	terminalConfig.AllowedCommands = append(terminalConfig.AllowedCommands, config.AllowedCommands{Command: "true", Args: []string{}})

	t.Run("command with no output", func(t *testing.T) {
		cmdString := "true"
		output, err := execute(cmdString, terminalConfig)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}

		expectedOutput := "No output"
		if output != expectedOutput {
			t.Errorf("expected output to be %q, but got %q", expectedOutput, output)
		}
	})
}
