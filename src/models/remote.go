package models

type RepositoryType uint8

const (
	Local RepositoryType = iota + 1
	Rclone
)

func RepositoryTypes() []SelectOption {
	var temp []SelectOption
	for i := 0; i < len(_RepositoryType_index)-1; i++ {
		temp = append(temp, SelectOption{
			Name:  _RepositoryType_name[_RepositoryType_index[i]:_RepositoryType_index[i+1]],
			Value: i + 1,
		})
	}
	return temp
}

type CompressionType uint8

const (
	Automatic CompressionType = iota + 1
	Maximum
	NoCompression
)

func CompressionTypes() []SelectOption {
	var temp []SelectOption
	for i := 0; i < len(_CompressionType_index)-1; i++ {
		temp = append(temp, SelectOption{
			Name:  _CompressionType_name[_CompressionType_index[i]:_CompressionType_index[i+1]],
			Value: i + 1,
		})
	}
	return temp
}

func ResticCompressionType(cType CompressionType) string {
	switch cType {
	case 1:
		return "auto"
	case 2:
		return "max"
	default:
		return "off"
	}
}

type Remote struct {
	ID              uint            `json:"id" form:"id" gorm:"primaryKey" validate:"omitempty,number"`
	Description     string          `json:"description" form:"description" validate:"required"`
	RepositoryType  RepositoryType  `json:"repository_type" form:"repository_type" validate:"required,oneof=1 2"`
	Repository      string          `json:"repository" form:"repository" validate:"required,endsnotwith=/"`
	PasswordFile    string          `json:"password_file" form:"password_file" validate:"required,file"`
	RetentionPolicy string          `json:"retention_policy" form:"retention_policy" validate:"-"`
	CompressionType CompressionType `json:"compression_type" form:"compression_type" validate:"required,oneof=1 2 3"`
	Jobs            []Job           `json:"-" form:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" validate:"-"`
}
