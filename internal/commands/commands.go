package commands

import (
	"os/exec"
)

func ExecuteCommand(cmdString string) (string, error) {
	cmd := exec.Command("sh", "-c", cmdString)
	out, err := cmd.CombinedOutput()
	result := string(out)
	if result == "" {
		result = "No output"
	}
	return result, err
}
