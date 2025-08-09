package services

import (
	"fmt"

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
	cs.Events.SendCommandEvent(fmt.Sprintf("Executing command: %s", cmdString))
	out, _ := commands.ExecuteCommand(cmdString)
	cs.Events.SendCommandEvent(out)
}
