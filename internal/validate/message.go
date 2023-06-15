package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func createInfoText(value interface{}, field string) string {
	return fmt.Sprintf("%s %s", value, field)
}

func MessageForError(err validator.FieldError) string {
	if strings.Contains(err.Namespace(), "Config") {
		switch err.Tag() {
		case "required", "rfg":
			return createInfoText(err.Field(), "is required")
		case "timezone":
			return createInfoText(err.Value(), "is not a valid Timezone")
		case "oneof":
			return createInfoText(err.Value(), "is no a valid option. Valid options are: "+err.Param())
		case "cron":
			return createInfoText(err.Value(), "is not a valid CRON expression")
		case "url":
			return createInfoText(err.Value(), "is not a valid URL")
		case "endswith":
			return createInfoText(err.Value(), "needs to have a trailing '/'")
		case "uuid":
			return createInfoText(err.Value(), "is not a valid UUID")
		case "shoutrrr":
			return createInfoText(err.Field(), "needs to be a valid Shoutrrr URL (https://containrrr.dev/shoutrrr)")
		default:
			return createInfoText(err.Value(), "is not valid")
		}
	} else {
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
		default:
			return createInfoText(err.Value(), "is not valid")
		}
	}
}
