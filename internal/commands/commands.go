package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

func PrepareCommand(command string) (program string, args []string, err error) {
	split := strings.Split(command, " ")

	if len(split) >= 2 {
		return split[0], split[1:], nil
	} else {
		return "", nil, fmt.Errorf("failed to parse command: %s", command)
	}
}

func ExecuteCommand(program string, args []string) (string, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
