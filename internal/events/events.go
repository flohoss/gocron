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
	EventStatus = "status"
)

type EventInfo struct {
	Idle bool `json:"idle"`
	Jobs any  `json:"jobs"`
	Runs any  `json:"runs"`
}

func New(onSubscribe func(streamID string, sub *sse.Subscriber)) *Event {
	sse := sse.NewWithCallback(onSubscribe, nil)
	sse.AutoReplay = false
	sse.CreateStream(EventStatus)
	return &Event{
		SSE: sse,
	}
}

func (e *Event) SendEvent(idle bool, jobs any, runs any) {
	data, _ := json.Marshal(&EventInfo{
		Idle: idle,
		Jobs: jobs,
		Runs: runs,
	})
	e.SSE.Publish(EventStatus, &sse.Event{
		Data: data,
	})
}

func (e *Event) GetHandler() echo.HandlerFunc {
	return echo.WrapHandler(http.HandlerFunc(e.SSE.ServeHTTP))
}
