package database

import (
	"encoding/json"

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

var TimeToGoBackInMilliseconds int64 = 604800000
