package models

type JobStatus uint8

const (
	New JobStatus = iota
	Running
	Stopped
	Success
)

type DatabaseType uint8

const (
	NoDatabase DatabaseType = iota + 1
	PostgreSQL
	MariaDB
)

func DatabaseTypes() []SelectOption {
	var temp []SelectOption
	for i := 0; i < len(_DatabaseType_index)-1; i++ {
		temp = append(temp, SelectOption{
			Name:  _DatabaseType_name[_DatabaseType_index[i]:_DatabaseType_index[i+1]],
			Value: i + 1,
		})
	}
	return temp
}

type Job struct {
	ID                uint         `json:"id" form:"id" gorm:"primaryKey" validate:"omitempty,number"`
	Description       string       `json:"description" form:"description" validate:"required"`
	LocalDirectory    string       `json:"local_directory" form:"local_directory" validate:"required,endsnotwith=/,dir"`
	RemoteDirectory   string       `json:"remote_directory" form:"remote_directory" validate:"required,startsnotwith=/,endsnotwith=/"`
	DatabaseType      DatabaseType `json:"database_type" form:"database_type" validate:"required,oneof=1 2 3"`
	DatabaseContainer string       `json:"database_container" form:"database_container" validate:"rfg=DatabaseType:1"`
	DatabaseUser      string       `json:"database_user" form:"database_user" validate:"rfg=DatabaseType:1"`
	DatabasePassword  string       `json:"database_password" form:"database_password" gorm:"-" validate:"rfg=DatabaseType:1"`
	DatabaseName      string       `json:"database_name" form:"database_name" validate:"rfg=DatabaseType:1"`
	RemoteID          uint         `json:"remote_id" form:"remote_id" validate:"required,number"`
	Remote            Remote       `json:"-" form:"-" validate:"-"`
	CustomCommand     string       `json:"custom_command" form:"custom_command" validate:"omitempty,ascii"`
	DockerRestart     bool         `json:"docker_restart" form:"docker_restart" validate:"omitempty"`
	CheckResticRepo   bool         `json:"check_restic_repo" form:"check_restic_repo" validate:"omitempty"`
	CreatedAt         int64        `json:"created_at" form:"-" validate:"-" gorm:"autoCreateTime"`
	UpdatedAt         int64        `json:"updated_at" form:"-" validate:"-" gorm:"autoUpdateTime"`
	JobStatus         JobStatus    `json:"job_status" form:"-" validate:"-"`
	Logs              []Log        `json:"-" form:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"-"`
}
