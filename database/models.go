// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package database

import (
	"database/sql"
	"time"
)

type Command struct {
	CommandID   int64
	JobID       int64
	CommandType string
	Command     string
}

type Job struct {
	JobID             int64
	Description       string
	LocalDirectory    string
	ResticRemote      string
	RestartOption     int64
	PasswordFilePath  string
	CompressionType   string
	SvgIcon           string
	CreatedAt         time.Time
	RetentionPolicyID int64
}

type Log struct {
	LogID   int64
	JobID   int64
	LogType string
	Message string
}

type RetentionPolicy struct {
	RetentionPolicyID int64
	RetentionPolicy   string
}

type Run struct {
	RunID     int64
	JobID     int64
	StartTime time.Time
	EndTime   time.Time
	LogID     sql.NullInt64
}
