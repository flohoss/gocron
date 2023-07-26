package database

import "github.com/labstack/echo/v4"

func (s *Service) CreateOrUpdateFromRequest(ctx echo.Context, value interface{}) error {
	if err := ctx.Validate(value); err != nil {
		return err
	}
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
