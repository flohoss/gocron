package services

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"gitlab.unjx.de/flohoss/gocron/config"
	"gitlab.unjx.de/flohoss/gocron/internal/commands"
	"gitlab.unjx.de/flohoss/gocron/internal/events"
)

type CommandsService struct {
	Events *events.Event
}

func NewCommandService(e *events.Event) *CommandsService {
	return &CommandsService{
		Events: e,
	}
}

func (cs *CommandsService) ExecuteCommand(cmdString string) {
	cs.Events.SendCommandEvent(int(Debug), fmt.Sprintf("Executing command: %s", cmdString))
	out, err := execute(cmdString, config.GetTerminalSettings())
	if err != nil {
		cs.Events.SendCommandEvent(int(Error), fmt.Sprintf("Error executing command: %s, %s", cmdString, err.Error()))
	}
	cs.Events.SendCommandEvent(int(Info), out)
}

func execute(cmdString string, settings config.TerminalSettings) (string, error) {
	if settings.AllowAllCommands {
		return commands.ExecuteCommand(cmdString)
	}

	parts := strings.Fields(cmdString)
	if len(parts) == 0 {
		return "", errors.New("empty command")
	}

	cmdName := parts[0]
	args := parts[1:]

	cmdConfig, found := settings.AllowedCommands[cmdName]
	if !found {
		return "", fmt.Errorf("command %q is not allowed", cmdName)
	}

	if len(cmdConfig.Args) > 0 {
		if len(args) > 0 {
			for _, arg := range args {
				if _, found := cmdConfig.AllowedArgsMap[arg]; !found {
					return "", fmt.Errorf("argument %q is not allowed for command %q", arg, cmdName)
				}
			}
		}
	} else if len(args) > 0 && !cmdConfig.AllowAllArgs {
		return "", fmt.Errorf("command %q does not allow any arguments", cmdName)
	}

	cmd := exec.Command(cmdName, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	if len(out) == 0 {
		return "No output", nil
	}
	return string(out), nil
}
