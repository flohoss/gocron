package controller

import (
	"encoding/json"
	"fmt"
	"gobackup/models"
	"time"

	"github.com/r3labs/sse/v2"
)

const (
	EventLog    string = "logs"
	EventSystem string = "logssystem"
	EventStatus string = "status"
)

type StatusMessage struct {
	JobID        uint             `json:"job_id"`
	Status       models.JobStatus `json:"status"`
	UpdateAt     string           `json:"updated_at"`
	JobsFinished bool             `json:"jobs_finished"`
}

func (c *Controller) setupEventChannel() {
	c.SSE.AutoReplay = false
	var jobs []models.Job
	c.orm.Find(&jobs)
	for _, job := range jobs {
		c.SSE.CreateStream(fmt.Sprintf("%s%d", EventLog, job.ID))
	}
	c.SSE.CreateStream(EventLog)
	c.SSE.CreateStream(EventSystem)
	c.SSE.CreateStream(EventStatus)
}

func (c *Controller) updateJobStatus(job *models.Job, status models.JobStatus) {
	job.JobStatus = status
	c.orm.Save(job)
	json, _ := json.Marshal(StatusMessage{JobID: job.ID, Status: status, UpdateAt: time.Unix(job.UpdatedAt, 0).Format("2006-01-02 15:04:05"), JobsFinished: false})
	c.SSE.Publish(EventStatus, &sse.Event{Data: json})
}

func (c *Controller) setJobsRunning(running bool) {
	c.JobRunning = running
	json, _ := json.Marshal(StatusMessage{JobID: 0, Status: models.New, UpdateAt: "", JobsFinished: true})
	c.SSE.Publish(EventStatus, &sse.Event{Data: json})
}
