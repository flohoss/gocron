package database

type RetentionPolicy struct {
	ID          uint64 `gorm:"primaryKey"`
	Description string `gorm:"unique"`
	Policy      string `gorm:"unique"`
	Jobs        []Job  `gorm:"constraint:OnDelete:SET NULL;"`
}

type CompressionType struct {
	ID          uint64 `gorm:"primaryKey"`
	Compression string `gorm:"unique"`
	Jobs        []Job  `gorm:"constraint:OnDelete:SET NULL;"`
}

type PreCommand struct {
	ID      uint64 `gorm:"primaryKey" form:"id" json:"id" validate:"omitempty"`
	JobId   uint64 `form:"job_id" json:"job_id" validate:"required"`
	Job     Job    `form:"-" json:"-" validate:"-"`
	Command string `form:"command" json:"command" validate:"required" example:"docker compose stop"`
}

type PostCommand struct {
	ID      uint64 `gorm:"primaryKey" form:"id" json:"id" validate:"omitempty"`
	JobId   uint64 `form:"job_id" json:"job_id" validate:"omitempty"`
	Job     Job    `form:"-" json:"-" validate:"-"`
	Command string `form:"command" json:"command" validate:"required" example:"docker compose up -d"`
}

type Log struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	RunID         uint64      `json:"run_id"`
	Run           Run         `json:"-"`
	LogTypeID     uint64      `json:"log_type_id"`
	LogType       LogType     `json:"-"`
	LogSeverityID uint64      `json:"log_severity_id"`
	LogSeverity   LogSeverity `json:"-"`
	Message       string      `json:"message"`
}

type LogType struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Type string `gorm:"unique" json:"type"`
	Logs []Log  `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
}

type LogSeverity struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Severity string `gorm:"unique" json:"severity"`
	Logs     []Log  `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
}

type Run struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	JobID     uint64 `json:"job_id"`
	Job       Job    `json:"-"`
	StartTime int64  `gorm:"autoCreateTime:milli" json:"start_time"`
	EndTime   int64  `json:"end_time"`
	LogID     uint64 `json:"log_id"`
	Logs      []Log  `json:"-" form:"-" gorm:"constraint:OnDelete:CASCADE;" validate:"-"`
}

type Job struct {
	ID                uint64          `gorm:"primaryKey" form:"id" json:"id" validate:"omitempty"`
	Description       string          `json:"description" validate:"required" example:"Gitea"`
	LocalDirectory    string          `json:"local_directory" validate:"required" example:"/opt/docker/gitea"`
	ResticRemote      string          `json:"restic_remote" validate:"required" example:"rclone:pcloud:Backups/gitea"`
	PasswordFilePath  string          `json:"password_file_path" validate:"required" example:"/secrets/.resticpwd"`
	SvgIcon           string          `json:"svg_icon" validate:"omitempty" example:""`
	RetentionPolicyID uint64          `json:"retention_policy_id" validate:"required,oneof=1 2 3 4 5 6 7 8 9" example:"1"`
	RetentionPolicy   RetentionPolicy `json:"-" validate:"-"`
	CompressionTypeID uint64          `form:"compression_type_id" json:"compression_type_id" validate:"required,oneof=1 2 3" example:"1"`
	CompressionType   CompressionType `json:"-" validate:"-"`
	PreCommands       []PreCommand    `json:"pre_commands" gorm:"constraint:OnDelete:CASCADE;" validate:"omitempty"`
	PostCommands      []PostCommand   `json:"post_commands" gorm:"constraint:OnDelete:CASCADE;" validate:"omitempty"`
	Runs              []Run           `json:"-" gorm:"constraint:OnDelete:CASCADE;" validate:"omitempty"`
}
