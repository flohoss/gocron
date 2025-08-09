package events

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
)

type Event struct {
	SSE *sse.Server
}

const (
	EventStatus  = "status"
	CommandEvent = "command"
)

type EventInfo struct {
	Idle bool `json:"idle"`
	Run  any  `json:"run"`
}

func New(onSubscribe func(streamID string, sub *sse.Subscriber)) *Event {
	sse := sse.NewWithCallback(onSubscribe, nil)
	sse.AutoReplay = false
	sse.CreateStream(EventStatus)
	sse.CreateStream(CommandEvent)
	return &Event{
		SSE: sse,
	}
}

func (e *Event) SendJobEvent(idle bool, run any) {
	data, _ := json.Marshal(&EventInfo{
		Idle: idle,
		Run:  run,
	})
	e.SSE.Publish(EventStatus, &sse.Event{
		Data: data,
	})
}

func (e *Event) SendCommandEvent(command string) {
	e.SSE.Publish(CommandEvent, &sse.Event{
		Data: []byte(command),
	})
}

func (e *Event) GetHandler() echo.HandlerFunc {
	return echo.WrapHandler(http.HandlerFunc(e.SSE.ServeHTTP))
}
