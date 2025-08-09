package services

import (
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

func (cs *CommandsService) ExecuteCommand(cmdString string) (string, error) {
	cs.Events.SendCommandEvent(cmdString)
	return commands.ExecuteCommand(cmdString)
}
