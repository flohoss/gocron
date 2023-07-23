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
