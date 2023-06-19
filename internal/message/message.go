package message

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func createInfoText(value interface{}, field string) string {
	return fmt.Sprintf("%s %s", value, field)
}

func MessageForError(err validator.FieldError) string {
	switch err.Tag() {
	case "required", "rfg":
		return createInfoText("This field", "is required")
	case "dir":
		return createInfoText("This", "is not an existing Directory. Please mount it first.")
	case "file":
		return createInfoText("This", "is not an existing file. Please mount if first")
	case "number":
		return createInfoText("This", "is not a valid number")
	case "oneof":
		return createInfoText("This", "is no a valid option. Valid options are: "+err.Param())
	case "endsnotwith":
		return createInfoText("This", "cannot end with '"+err.Param()+"'")
	case "startsnotwith":
		return createInfoText("This", "cannot start with '"+err.Param()+"'")
	case "ascii":
		return createInfoText("This", "can only contain ASCII character")
	default:
		return createInfoText(err.Value(), "is not valid")
	}
}
