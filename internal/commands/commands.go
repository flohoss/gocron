package commands

import (
	"database/sql"
	"os"
	"os/exec"
	"strings"
)

func ExtractVariable(content string) string {
	var result strings.Builder
	start := 0

	for {
		// Find the start of a variable
		startIdx := strings.Index(content[start:], "${")
		if startIdx == -1 {
			// No more variables, append the remaining content and break
			result.WriteString(content[start:])
			break
		}
		startIdx += start // Adjust relative index to absolute

		// Find the end of the variable
		endIdx := strings.IndexByte(content[startIdx:], '}')
		if endIdx == -1 {
			// No closing brace, append the remaining content and break
			result.WriteString(content[start:])
			break
		}
		endIdx += startIdx // Adjust relative index to absolute

		// Append the part before the variable
		result.WriteString(content[start:startIdx])

		// Extract and resolve the variable
		varName := content[startIdx+2 : endIdx]
		envValue := os.Getenv(varName)
		if envValue != "" {
			result.WriteString(envValue) // Append the resolved value
		} else {
			result.WriteString(content[startIdx : endIdx+1]) // Keep the unresolved variable
		}

		// Move the start position after the variable
		start = endIdx + 1
	}

	return result.String()
}

func PrepareCommand(command string) (program string, args []string) {
	split := strings.Fields(command)
	if len(split) == 0 {
		return "", nil
	}

	for i := 0; i < len(split); {
		expanded := ExtractVariable(split[i]) // Expand the variable
		envParts := strings.Fields(expanded)  // Split expanded value into parts

		// Replace the current element with the expanded parts
		split = append(split[:i], append(envParts, split[i+1:]...)...)

		// Move the index forward to skip over the newly added parts
		i += len(envParts)
	}

	if len(split) > 0 {
		return split[0], split[1:]
	}

	return "", nil
}

func ExecuteCommand(program string, args []string, fileOutput sql.NullString) (string, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	if fileOutput.Valid {
		file, err := os.OpenFile(ExtractVariable(fileOutput.String), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err.Error(), err
		}
		defer file.Close()
		file.Write(out)
		return fileOutput.String, nil
	}
	return string(out), err
}
