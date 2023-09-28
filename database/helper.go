package database

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/internal/message"
)

func (s *Service) ValidateRequestBinding(ctx echo.Context, value interface{}) ([]byte, error) {
	if err := ctx.Validate(value); err != nil {
		tmp := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			tmp[err.Field()] = message.MessageForError(err)
		}
		jsonData, _ := json.Marshal(tmp)
		return jsonData, err
	}
	return []byte{}, nil
}

func (s *Service) CreateOrUpdateFromRequest(ctx echo.Context, value interface{}) error {
	if err := s.orm.Save(value).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateOrUpdate(value interface{}) error {
	if logValue, ok := value.(*Log); ok {
		if logValue.LogSeverity == LogWarning {
			s.notifyService.SendNotification(fmt.Sprintf("%s Warning - %s", emoji.Warning, s.identifier), logValue.Message)
		} else if logValue.LogSeverity == LogError {
			s.notifyService.SendNotification(fmt.Sprintf("%s Error - %s", emoji.CrossMark, s.identifier), logValue.Message)
		}
	}
	if err := s.orm.Save(value).Error; err != nil {
		return err
	}
	return nil
}

func IsInArray(arr []Command, target Command) bool {
	for _, item := range arr {
		if item.ID == target.ID {
			return true
		}
	}
	return false
}

func AnonymisePasswords(text string) string {
	pattern := `((?:--password|PASSWORD|-p)[="'\s]+)(.+?)(["'\s])`
	re := regexp.MustCompile(pattern)
	anonymizedText := strings.TrimSpace(re.ReplaceAllString(text+" ", `${1}****${3}`))
	return anonymizedText
}
