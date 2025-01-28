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
	Job  Job  `json:"job"`
}

type Job struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Log  Log    `json:"log"`
}

type Log struct {
	CreatedAt  int64  `json:"created_at"`
	SeverityID int64  `json:"severity_id"`
	Message    string `json:"message"`
}

func New(jobs []string) *Event {
	sse := sse.New()
	sse.AutoReplay = false
	sse.CreateStream(EventStatus)
	return &Event{
		SSE: sse,
	}
}

func (e *Event) SendEvent(info *EventInfo) {
	data, _ := json.Marshal(info)
	e.SSE.Publish(EventStatus, &sse.Event{
		Data: data,
	})
}

func (e *Event) GetHandler() echo.HandlerFunc {
	return echo.WrapHandler(http.HandlerFunc(e.SSE.ServeHTTP))
}
