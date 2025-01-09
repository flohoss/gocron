package commands

import (
	"database/sql"
	"os"
	"os/exec"
	"strings"
)

func PrepareCommand(command string) (program string, args []string) {
	split := strings.Split(command, " ")

	if len(split) >= 2 {
		for i := 1; i < len(split); i++ {
			if strings.HasPrefix(split[i], "$") {
				envValue := os.Getenv(split[i][1:])
				if envValue != "" {
					envParts := strings.Split(envValue, " ")
					split = append(split[:i], append(envParts, split[i+1:]...)...)
					i += len(envParts) - 1
				} else {
					split = append(split[:i], split[i+1:]...)
				}
			} else {
				i++
			}
		}
		return split[0], split[1:]
	} else {
		return split[0], []string{}
	}
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
