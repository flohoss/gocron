package events

import (
	"database/sql"
	"encoding/json"

	"github.com/r3labs/sse/v2"
)

type Event struct {
	SSE *sse.Server
}

const (
	EventGlobal = "global"
	EventJob    = "job_"
)

type GlobalInfo struct {
	Idle bool   `json:"idle"`
	Jobs []Jobs `json:"jobs"`
}

type Jobs struct {
	Name string      `json:"name"`
	Cron string      `json:"cron"`
	Runs interface{} `json:"runs"`
}

type Runs struct {
	JobID    string         `json:"job_id"`
	StatusID int64          `json:"status_id"`
	Duration sql.NullString `json:"duration"`
}

type JobInfo struct {
	Log Log `json:"log"`
}

type Log struct {
	CreatedAt  int64  `json:"created_at"`
	SeverityID int64  `json:"severity_id"`
	Message    string `json:"message"`
}

func New(jobs []string) *Event {
	sse := sse.New()
	sse.CreateStream(EventGlobal)
	for _, job := range jobs {
		sse.CreateStream(EventJob + job)
	}
	return &Event{
		SSE: sse,
	}
}

func (e *Event) SendGlobal(info *GlobalInfo) {
	data, _ := json.Marshal(info)
	e.SSE.Publish(EventGlobal, &sse.Event{
		Data: data,
	})
}

func (e *Event) SendJob(jobID string, info *JobInfo) {
	data, _ := json.Marshal(info)
	e.SSE.Publish(EventJob+jobID, &sse.Event{
		Data: data,
	})
}
