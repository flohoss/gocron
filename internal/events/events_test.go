package events

import (
	"testing"

	"github.com/r3labs/sse/v2"
)

func TestNew_CreatesStreams(t *testing.T) {
	e := New(func(streamID string, sub *sse.Subscriber) {})

	if e.SSE == nil {
		t.Fatal("expected SSE server to be initialized")
	}
	if !e.SSE.StreamExists(EventStatus) {
		t.Fatalf("expected stream %q to exist", EventStatus)
	}
	if !e.SSE.StreamExists(CommandEvent) {
		t.Fatalf("expected stream %q to exist", CommandEvent)
	}
}

func TestSendJobEvent_DoesNotPanic(t *testing.T) {
	e := New(func(streamID string, sub *sse.Subscriber) {})
	// Publish to a stream with no subscribers must not panic
	e.SendJobEvent(true, "run-1", []string{"job-a"})
	e.SendJobEvent(false, nil, nil)
}

func TestSendCommandEvent_DoesNotPanic(t *testing.T) {
	e := New(func(streamID string, sub *sse.Subscriber) {})
	e.SendCommandEvent(2, "echo hello")
	e.SendCommandEvent(0, "")
}

func TestGetHandler_ReturnsNonNilHandler(t *testing.T) {
	e := New(func(streamID string, sub *sse.Subscriber) {})
	if e.GetHandler() == nil {
		t.Fatal("expected non-nil handler")
	}
}
