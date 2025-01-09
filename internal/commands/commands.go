package commands

import (
	"database/sql"
	"os"
	"os/exec"
	"strings"
)

func PrepareCommand(command string) (program string, args []string) {
	split := strings.Fields(command)
	if len(split) == 0 {
		return "", nil
	}

	for i := 1; i < len(split); {
		if strings.HasPrefix(split[i], "$") {
			envValue := os.Getenv(split[i][1:])
			if envValue != "" {
				envParts := strings.Fields(envValue)
				split = append(split[:i], append(envParts, split[i+1:]...)...)
				i += len(envParts)
			} else {
				split = append(split[:i], split[i+1:]...)
			}
		} else {
			i++
		}
	}

	return split[0], split[1:]
}

func ExecuteCommand(program string, args []string, fileOutput sql.NullString) (string, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	if fileOutput.Valid {
		file, err := os.OpenFile(fileOutput.String, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return "Could not create file", err
		}
		defer file.Close()
		file.Write(out)
		return fileOutput.String, nil
	}
	return string(out), err
}
