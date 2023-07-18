package controller

import (
	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"go.uber.org/zap"
)

type Controller struct {
	service *database.Service
	env     *env.Config
}

type IndexData struct {
	Title string
	Jobs  []database.Job
}

func NewController(env *env.Config) *Controller {
	service, err := database.MigrateDatabase()
	if err != nil {
		zap.L().Fatal("cannot connect to database", zap.Error(err))
	}
	ctrl := Controller{service: service, env: env}
	return &ctrl
}
