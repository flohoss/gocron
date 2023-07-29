package database

type CommandInfo struct {
	Description string
	Command     string
}

type RetentionPolicy uint8

const (
	KeepAll RetentionPolicy = iota + 1
	KeepDailyLast2
	KeepDailyLast7
	KeepDailyLast31
	KeepMostRecent7Daily
	KeepMostRecent31Daily
	KeepDailyFor5Years
)

var RetentionPolicyInfoMap = map[RetentionPolicy]CommandInfo{
	KeepAll:               {"Keep all snapshots", ""},
	KeepDailyLast2:        {"Keep daily snapshots for the last 2 days", "--keep-daily 2"},
	KeepDailyLast7:        {"Keep daily snapshots for the last 7 days", "--keep-daily 7"},
	KeepDailyLast31:       {"Keep daily snapshots for the last 31 days", "--keep-daily 31"},
	KeepMostRecent7Daily:  {"Keep the most recent 7 daily, 4 last-day-of-the-week, 12 or 11 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots", "--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 11"},
	KeepMostRecent31Daily: {"Keep the most recent 31 daily, 8 last-day-of-the-week, 24 or 23 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots", "--keep-daily 31 --keep-weekly 9 --keep-monthly 24 --keep-yearly 11"},
	KeepDailyFor5Years:    {"Keep daily for 5 Years, 520 last-day-of-the-week, 121 or 120 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots", "--keep-daily 1095 --keep-weekly 521 --keep-monthly 121 --keep-yearly 11"},
}

type CompressionType uint8

const (
	Automatic CompressionType = iota + 1
	Maximum
	NoCompression
)

var CompressionTypeInfoMap = map[CompressionType]CommandInfo{
	Automatic:     {"Automatic", "auto"},
	Maximum:       {"Maximum", "max"},
	NoCompression: {"No compression", "off"},
}

type SystemLog struct {
	ID          uint64      `gorm:"primaryKey" json:"id"`
	Message     string      `json:"message"`
	CreatedAt   int64       `gorm:"autoCreateTime:milli" json:"created_at"`
	LogSeverity LogSeverity `json:"log_severity"`
}

type LogType uint8

const (
	LogGeneral LogType = iota + 1
	LogRestic
	LogCustom
	LogPrune
	LogCheck
)

type LogSeverity uint8

const (
	LogNone LogSeverity = iota
	LogInfo
	LogWarning
	LogError
	LogRunning
)

type Log struct {
	ID          uint64      `gorm:"primaryKey" json:"id" validate:"omitempty"`
	RunID       uint64      `json:"run_id" validate:"omitempty"`
	Run         Run         `json:"-" validate:"-"`
	LogType     LogType     `json:"log_type" validate:"required"`
	LogSeverity LogSeverity `json:"log_severity" validate:"required"`
	Message     string      `json:"message" validate:"required"`
	CreatedAt   int64       `gorm:"autoCreateTime:milli" json:"created_at"`
}

type Run struct {
	ID        uint64 `gorm:"primaryKey" json:"id" validate:"omitempty"`
	JobID     uint64 `json:"job_id" validate:"omitempty"`
	Job       Job    `json:"-" validate:"-"`
	StartTime int64  `gorm:"autoCreateTime:milli" json:"start_time"`
	EndTime   int64  `json:"end_time" validate:"omitempty"`
	Logs      []Log  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"logs" validate:"omitempty"`
}

type Command struct {
	ID         uint64 `gorm:"primaryKey" json:"id" validate:"omitempty"`
	SortID     uint64 `json:"sort_id" validate:"required"`
	Type       uint8  `json:"type" validate:"required,oneof=1 2"`
	JobId      uint64 `json:"job_id" validate:"omitempty"`
	Job        Job    `json:"-" validate:"-"`
	Command    string `json:"command" validate:"required,ascii" example:"docker compose stop"`
	FileOutput string `json:"file_output" validate:"omitempty" example:".dbBackup.sql"`
}

type Job struct {
	ID               uint64          `gorm:"primaryKey" json:"id" validate:"omitempty"`
	Description      string          `json:"description" validate:"required" example:"Gitea"`
	LocalDirectory   string          `json:"local_directory" validate:"required,dir" example:"/opt/docker/gitea"`
	ResticRemote     string          `json:"restic_remote" validate:"required" example:"rclone:pcloud:Backups/gitea"`
	PasswordFilePath string          `json:"password_file_path" validate:"required,file" example:"/secrets/.resticpwd"`
	RetentionPolicy  RetentionPolicy `json:"retention_policy" validate:"required,oneof=1 2 3 4 5 6 7" example:"1"`
	CompressionType  CompressionType `json:"compression_type" validate:"required,oneof=1 2 3" example:"1"`
	RoutineCheck     uint64          `json:"routine_check" validate:"omitempty,number,min=0,max=100"`
	PreCommands      []Command       `json:"pre_commands" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"omitempty"`
	PostCommands     []Command       `json:"post_commands" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"omitempty"`
	Runs             []Run           `json:"runs" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"omitempty"`
	Status           LogSeverity     `json:"status" gorm:"-" validate:"-"`
}
