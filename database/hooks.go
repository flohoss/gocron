package database

import (
	"encoding/json"

	"github.com/r3labs/sse/v2"
	"gorm.io/gorm"
)

var SSE *sse.Server

func SetupEventChannel() {
	SSE = sse.New()
	SSE.AutoReplay = false
	SSE.CreateStream("runs")
	SSE.CreateStream("logs")
}

func (r *Run) AfterCreate(tx *gorm.DB) (err error) {
	json, _ := json.Marshal(r)
	SSE.Publish("runs", &sse.Event{Data: json})
	return
}

func (r *Run) AfterUpdate(tx *gorm.DB) (err error) {
	run := new(Run)
	tx.Preload("Logs").Find(run, r.ID)
	json, _ := json.Marshal(run)
	SSE.Publish("runs", &sse.Event{Data: json})
	return
}

func (l *Log) AfterCreate(tx *gorm.DB) (err error) {
	run := new(Run)
	tx.Find(run, l.RunID)
	run.Logs = append(run.Logs, *l)
	json, _ := json.Marshal(run)
	SSE.Publish("logs", &sse.Event{Data: json})
	return
}

func (l *Log) AfterUpgrade(tx *gorm.DB) (err error) {
	run := new(Run)
	tx.Find(run, l.RunID)
	run.Logs = append(run.Logs, *l)
	json, _ := json.Marshal(run)
	SSE.Publish("logs", &sse.Event{Data: json})
	return
}
