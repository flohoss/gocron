// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package jobs

import (
	"database/sql"
	"time"
)

type BackupSchedule struct {
	ID             int64          `json:"id"`
	CronExpression string         `json:"cron_expression"`
	Description    sql.NullString `json:"description"`
}

type Command struct {
	ID           int64  `json:"id"`
	JobID        int64  `json:"job_id"`
	CommandOrder int64  `json:"command_order"`
	CommandText  string `json:"command_text"`
}

type Env struct {
	ID    int64  `json:"id"`
	JobID int64  `json:"job_id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Job struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	BackupScheduleID int64  `json:"backup_schedule_id"`
}

type Log struct {
	ID          int64     `json:"id"`
	RunID       int64     `json:"run_id"`
	LogSeverity string    `json:"log_severity"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}

type Run struct {
	ID        int64        `json:"id"`
	JobID     int64        `json:"job_id"`
	StartTime time.Time    `json:"start_time"`
	EndTime   sql.NullTime `json:"end_time"`
	Status    string       `json:"status"`
}
