package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/controller"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"gitlab.unjx.de/flohoss/gobackup/internal/logging"
	"gitlab.unjx.de/flohoss/gobackup/internal/router"
	"go.uber.org/zap"
)

func main() {
	env, err := env.Parse()
	if err != nil {
		log.Fatal(err)
	}
	zap.ReplaceGlobals(logging.CreateLogger(env.LogLevel))

	queries, err := database.MigrateDatabase()
	if err != nil {
		zap.L().Error(err.Error())
	}
	ctx := context.Background()
	jobs, err := queries.ListJobs(ctx)
	if err != nil {
		zap.L().Error(err.Error())
	}
	zap.L().Info("query finished", zap.Any("jobs", jobs))

	r := router.InitRouter()
	c := controller.NewController(env)
	router.SetupRoutes(r, c)

	zap.L().Info("starting server", zap.String("url", fmt.Sprintf("http://localhost:%d", env.Port)), zap.String("version", env.Version))
	if err := r.Start(fmt.Sprintf(":%d", env.Port)); err != http.ErrServerClosed {
		zap.L().Fatal("cannot start server", zap.Error(err))
	}
}
