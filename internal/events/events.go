package events

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

type Event struct {
	SSE *sse.Server
}

const (
	EventStatus = "status"
)

type EventInfo struct {
	Idle bool           `json:"idle"`
	Data *jobs.JobsView `json:"data"`
}

func New(jobs []string, onSubscribe func(streamID string, sub *sse.Subscriber)) *Event {
	sse := sse.NewWithCallback(onSubscribe, nil)
	sse.AutoReplay = false
	sse.CreateStream(EventStatus)
	return &Event{
		SSE: sse,
	}
}

func (e *Event) SendEvent(idle bool, info *jobs.JobsView) {
	data, _ := json.Marshal(&EventInfo{
		Idle: idle,
		Data: info,
	})
	e.SSE.Publish(EventStatus, &sse.Event{
		Data: data,
	})
}

func (e *Event) GetHandler() echo.HandlerFunc {
	return echo.WrapHandler(http.HandlerFunc(e.SSE.ServeHTTP))
}
