package handlers

import (
	"net/http"
	"testing"
)

type mockCommandsService struct{}

func (m *mockCommandsService) ExecuteCommand(cmdString string) {}

func TestNewCommandHandler_AssignsService(t *testing.T) {
	mock := &mockCommandsService{}
	h := NewCommandHandler(mock)

	if h.CommandsService != mock {
		t.Fatal("expected command service to be assigned")
	}
}

func TestExecuteCommandOperation_HasExpectedMetadata(t *testing.T) {
	h := &CommandHandler{}
	op := h.executeCommandOperation()

	if op.OperationID != "post-command" {
		t.Fatalf("unexpected operation id: %q", op.OperationID)
	}
	if op.Method != http.MethodPost {
		t.Fatalf("unexpected method: %q", op.Method)
	}
	if op.Path != "/api/command" {
		t.Fatalf("unexpected path: %q", op.Path)
	}
}
