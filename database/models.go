package database

type RetentionPolicy struct {
	ID          uint64 `gorm:"primaryKey"`
	Description string `gorm:"unique"`
	Policy      string `gorm:"unique"`
	Jobs        []Job  `gorm:"constraint:OnDelete:SET NULL;"`
}

type CompressionType struct {
	ID          uint64 `gorm:"primaryKey"`
	Description string `gorm:"unique"`
	Compression string `gorm:"unique"`
	Jobs        []Job  `gorm:"constraint:OnDelete:SET NULL;"`
}

type LogType struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Type string `gorm:"unique" json:"type"`
	Logs []Log  `gorm:"constraint:OnDelete:SET NULL;" json:"logs" validate:"-"`
}

type LogSeverity struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Severity string `gorm:"unique" json:"severity"`
	Logs     []Log  `gorm:"constraint:OnDelete:SET NULL;" json:"logs" validate:"-"`
}

type SystemLog struct {
	ID            uint64 `gorm:"primaryKey" json:"id"`
	Message       string `json:"message"`
	CreatedAt     int64  `gorm:"autoCreateTime:milli" json:"created_at"`
	LogSeverityID uint64 `json:"log_severity_id"`
}

type Log struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	RunID         uint64      `json:"run_id"`
	LogTypeID     uint64      `json:"log_type_id"`
	LogType       LogType     `json:"-"`
	LogSeverityID uint64      `json:"log_severity_id"`
	LogSeverity   LogSeverity `json:"-"`
	Message       string      `json:"message"`
	CreatedAt     int64       `gorm:"autoCreateTime:milli" json:"created_at"`
}

type Run struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	JobID     uint64 `json:"job_id"`
	StartTime int64  `gorm:"autoCreateTime:milli" json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Logs      []Log  `gorm:"constraint:OnDelete:CASCADE;" json:"logs" validate:"-"`
}

type Command struct {
	ID      uint64 `gorm:"primaryKey" json:"id" validate:"omitempty"`
	SortID  uint64 `json:"sort_id" validate:"required"`
	Type    uint8  `json:"type" validate:"required,oneof=1 2"`
	JobId   uint64 `json:"job_id" validate:"omitempty"`
	Command string `json:"command" validate:"required,ascii" example:"docker compose stop"`
}

type Job struct {
	ID                uint64          `gorm:"primaryKey" json:"id" validate:"omitempty"`
	Description       string          `json:"description" validate:"required" example:"Gitea"`
	LocalDirectory    string          `json:"local_directory" validate:"required,dir" example:"/opt/docker/gitea"`
	ResticRemote      string          `json:"restic_remote" validate:"required" example:"rclone:pcloud:Backups/gitea"`
	PasswordFilePath  string          `json:"password_file_path" validate:"required,file" example:"/secrets/.resticpwd"`
	SvgIcon           string          `json:"svg_icon" validate:"omitempty" example:""`
	RetentionPolicyID uint64          `json:"retention_policy_id" validate:"required,oneof=1 2 3 4 5 6 7 8 9" example:"1"`
	RetentionPolicy   RetentionPolicy `json:"-" validate:"-"`
	CompressionTypeID uint64          `json:"compression_type_id" validate:"required,oneof=1 2 3" example:"1"`
	CompressionType   CompressionType `json:"-" validate:"-"`
	RoutineCheck      uint64          `json:"routine_check" validate:"omitempty,number,min=0,max=100"`
	PreCommands       []Command       `json:"pre_commands" gorm:"constraint:OnDelete:CASCADE;" validate:"omitempty"`
	PostCommands      []Command       `json:"post_commands" gorm:"constraint:OnDelete:CASCADE;" validate:"omitempty"`
	Runs              []Run           `json:"runs" gorm:"constraint:OnDelete:CASCADE;" validate:"-"`
}

type JobStats struct {
	TotalRuns   uint64 `json:"total_runs" validate:"required"`
	TotalLogs   uint64 `json:"total_logs" validate:"required"`
	InfoLogs    uint64 `json:"info_logs" validate:"required"`
	WarningLogs uint64 `json:"warning_logs" validate:"required"`
	ErrorLogs   uint64 `json:"error_logs" validate:"required"`
	ResticRuns  uint64 `json:"restic_runs" validate:"required"`
	CustomRuns  uint64 `json:"custom_runs" validate:"required"`
	PruneRuns   uint64 `json:"prune_runs" validate:"required"`
	CheckRuns   uint64 `json:"check_runs" validate:"required"`
}
