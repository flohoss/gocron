package main

import (
	"gobackup/controller"
	"gobackup/env"
	"gobackup/validate"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	val := validate.NewValidator()
	env := env.Parse(val)
	zap.ReplaceGlobals(createLogger(env.LogLevel))

	controller := controller.NewController(env)
	router := initRouter(val)
	setupRoutes(router, controller)

	go func() {
		zap.L().Info("starting server", zap.String("url", fmt.Sprintf("http://localhost:%d", env.Port)), zap.String("version", env.Version))
		if err := router.Start(fmt.Sprintf(":%d", env.Port)); err != http.ErrServerClosed {
			zap.L().Fatal("cannot start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	// https://docs.docker.com/engine/reference/commandline/stop/
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	controller.SSE.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		zap.L().Fatal("cannot shutdown server", zap.Error(err))
	}
}
