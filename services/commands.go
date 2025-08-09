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
	severity := Debug
	cs.Events.SendCommandEvent(int(severity), fmt.Sprintf("Executing command: %s", cmdString))
	out, err := commands.ExecuteCommand(cmdString)
	severity = Info
	if err != nil {
		severity = Error
	}
	cs.Events.SendCommandEvent(int(severity), out)
}
