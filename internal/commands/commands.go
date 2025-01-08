package commands

import (
	"os/exec"
	"strings"
)

func PrepareCommand(command string) (program string, args []string) {
	split := strings.Split(command, " ")

	if len(split) >= 2 {
		return split[0], split[1:]
	} else {
		return split[0], []string{}
	}
}

func ExecuteCommand(program string, args []string) (string, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
