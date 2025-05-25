package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"gitlab.unjx.de/flohoss/gocron/config"
	"gitlab.unjx.de/flohoss/gocron/handlers"
	"gitlab.unjx.de/flohoss/gocron/internal/env"
	"gitlab.unjx.de/flohoss/gocron/internal/notify"
	"gitlab.unjx.de/flohoss/gocron/internal/scheduler"
	"gitlab.unjx.de/flohoss/gocron/services"
)

const (
	configFolder = "./config/"
)

func init() {
	os.Mkdir(configFolder, os.ModePerm)
}

func setupRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "events")
		},
	}))

	return e
}

func main() {
	e := setupRouter()

	env, err := env.Parse()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	e.Logger.SetLevel(env.GetLogLevel())
	if env.GetLogLevel() == log.DEBUG {
		e.Use(middleware.Logger())
		e.Debug = true
	}

	cfg, err := config.New(configFolder + "config.yml")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	c := scheduler.New(env.DeleteRunsAfterDays)
	n := notify.New(env.AppriseUrl, env.GetAppriseNotifyLevel())

	js, err := services.NewJobService(configFolder+"db.sqlite", cfg, c, n)
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
	jh := handlers.NewJobHandler(js, n)

	handlers.SetupRouter(e, jh)

	e.Logger.Infof("Server starting on http://localhost:%d", env.Port)
	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	// Listen for OS signals to gracefully shut down
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", env.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt or SIGTERM signal to gracefully shut down the server.
	<-ctx.Done()
	e.Logger.Info("Received shutdown signal. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("Server shut down gracefully.")
}
