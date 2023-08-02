package message

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
)

func TestMessageForError(t *testing.T) {
	validate := env.NewEnvValidator()

	type MessageTest struct {
		Description      string `validate:"required"`
		LocalDirectory   string `validate:"required,dir"`
		PasswordFilePath string `validate:"required,file"`
		RetentionPolicy  uint8  `validate:"required,oneof=1 2 3 4 5 6 7"`
		CompressionType  uint8  `validate:"required,oneof=1 2 3"`
		RoutineCheck     uint64 `validate:"omitempty,number,min=0,max=100"`
	}

	msg := MessageTest{
		Description:      "",
		LocalDirectory:   "",
		PasswordFilePath: "",
		RetentionPolicy:  0,
		CompressionType:  0,
		RoutineCheck:     0,
	}
	err := validate.Struct(msg)
	for _, err := range err.(validator.ValidationErrors) {
		res := MessageForError(err)
		assert.NotContains(t, res, "is not valid")
	}

	msg = MessageTest{
		Description:      "example",
		LocalDirectory:   "/",
		PasswordFilePath: "/",
		RetentionPolicy:  1,
		CompressionType:  2,
		RoutineCheck:     90,
	}
	err = validate.Struct(msg)
	for _, err := range err.(validator.ValidationErrors) {
		res := MessageForError(err)
		assert.NotContains(t, res, "is not valid")
	}
}
