package database

import (
	"encoding/json"

	"github.com/r3labs/sse/v2"
	"gorm.io/gorm"
)

var SSE *sse.Server

type EventType uint8

const (
	EventCreateRun EventType = iota + 1
	EventUpdateRun
	EventCreateLog
)

type SSEvent struct {
	EventType EventType   `json:"event_type"`
	Content   interface{} `json:"content"`
}

func SetupEventChannel() {
	SSE = sse.New()
	SSE.AutoReplay = false
	SSE.CreateStream("jobs")
	SSE.CreateStream("restore_logs")
}

func (r *Run) AfterCreate(tx *gorm.DB) (err error) {
	json, _ := json.Marshal(SSEvent{EventType: EventCreateRun, Content: r})
	SSE.Publish("jobs", &sse.Event{Data: json})
	return
}

func (r *Run) AfterUpdate(tx *gorm.DB) (err error) {
	r.Status = r.getHighestLogSeverity(tx)
	json, _ := json.Marshal(SSEvent{EventType: EventUpdateRun, Content: r})
	SSE.Publish("jobs", &sse.Event{Data: json})
	return
}

func (l *Log) AfterCreate(tx *gorm.DB) (err error) {
	json, _ := json.Marshal(SSEvent{EventType: EventCreateLog, Content: l})
	SSE.Publish("jobs", &sse.Event{Data: json})
	return
}

func (l *SystemLog) AfterCreate(tx *gorm.DB) (err error) {
	json, _ := json.Marshal(l)
	SSE.Publish("restore_logs", &sse.Event{Data: json})
	return
}
