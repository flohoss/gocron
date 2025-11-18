package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/handlers"
	"github.com/flohoss/gocron/internal/events"
	"github.com/flohoss/gocron/internal/software"
	"github.com/flohoss/gocron/services"
)

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
	config.New()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.GetLogLevel(),
	}))
	slog.SetDefault(logger)

	software.Install()

	js, err := services.NewJobService()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	js.SetEvents(events.New(func(streamID string, sub *sse.Subscriber) {
		if streamID == events.EventStatus {
			js.Events.SendJobEvent(js.IsIdle(), nil, nil)
		}
	}))
	jh := handlers.NewJobHandler(js)

	cs := services.NewCommandService(js.Events)
	ch := handlers.NewCommandHandler(cs)

	handlers.SetupRouter(e, jh, ch)

	slog.Info("Starting server", "url", fmt.Sprintf("http://%s", config.GetServer()))
	slog.Error(e.Start(config.GetServer()).Error())
}
