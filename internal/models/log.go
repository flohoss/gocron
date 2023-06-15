package models

import "github.com/enescakir/emoji"

type LogTopic uint8

const (
	Backup LogTopic = iota + 1
	Restic
	Docker
	Database
)

type LogType uint8

const (
	Info LogType = iota + 1
	Warn
	Error
)

func (l LogType) Emoji() emoji.Emoji {
	switch l {
	case Warn:
		return emoji.Warning
	case Error:
		return emoji.CrossMark
	default:
		return emoji.Information
	}
}

func LogTypes() []SelectOption {
	var temp []SelectOption
	for i := 0; i < len(_LogType_index)-1; i++ {
		temp = append(temp, SelectOption{
			Name:  _LogType_name[_LogType_index[i]:_LogType_index[i+1]],
			Value: i + 1,
		})
	}
	return temp
}

type Log struct {
	CreatedAt int64    `json:"created_at" gorm:"primaryKey;autoIncrement:false;autoCreateTime:milli"`
	JobID     uint     `json:"job_id"`
	Job       Job      `json:"job"`
	Type      LogType  `json:"type"`
	Topic     LogTopic `json:"topic"`
	Message   string   `json:"message"`
}

type SystemLog struct {
	CreatedAt int64    `json:"created_at" gorm:"primaryKey;autoIncrement:false;autoCreateTime:milli"`
	Type      LogType  `json:"type"`
	Topic     LogTopic `json:"topic"`
	Message   string   `json:"message"`
}

type LogMessage struct {
	Description string   `json:"description"`
	Type        LogType  `json:"type"`
	Topic       LogTopic `json:"topic"`
	Message     string   `json:"message"`
	CreatedAt   int64    `json:"created_at"`
}
