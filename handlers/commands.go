package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type CommandsService interface {
	ExecuteCommand(cmdString string) (string, error)
}

func NewCommandHandler(cs CommandsService) *CommandHandler {
	return &CommandHandler{
		CommandsService: cs,
	}
}

type CommandHandler struct {
	CommandsService CommandsService
}

func (ch *CommandHandler) executeCommandOperation() huma.Operation {
	return huma.Operation{
		OperationID: "post-command",
		Method:      http.MethodPost,
		Path:        "/api/command",
		Summary:     "Execute command",
		Description: "Execute command.",
		Tags:        []string{"Command"},
	}
}

func (ch *CommandHandler) executeCommandHandler(ctx context.Context, input *struct {
	Command string `path:"command" minLength:"1" maxLength:"255" doc:"command to execute"`
}) (*struct{}, error) {
	go ch.CommandsService.ExecuteCommand(input.Command)
	return nil, nil
}
