package main

import (
	"fmt"
	"log"
	"net/http"

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

	r := router.InitRouter()
	c := controller.NewController(env)
	router.SetupRoutes(r, c)

	zap.L().Info("starting server", zap.String("url", fmt.Sprintf("http://localhost:%d", env.Port)), zap.String("version", env.Version))
	if err := r.Start(fmt.Sprintf(":%d", env.Port)); err != http.ErrServerClosed {
		zap.L().Fatal("cannot start server", zap.Error(err))
	}
}
