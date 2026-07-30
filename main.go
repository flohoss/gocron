package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/r3labs/sse/v2"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/handlers"
	"github.com/flohoss/gocron/internal/buildinfo"
	"github.com/flohoss/gocron/internal/cli"
	"github.com/flohoss/gocron/internal/events"
	"github.com/flohoss/gocron/internal/software"
	"github.com/flohoss/gocron/services"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	if opts.ShowVersion {
		fmt.Println(buildinfo.Summary())
		return
	}

	config.New(opts.ConfigFile)

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

	e := handlers.InitRouter()
	handlers.SetupRouter(e, jh, ch)

	slog.Info("Starting server", "url", fmt.Sprintf("http://%s", config.GetServer()))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down scheduler and running jobs")
		js.Shutdown()
	}()

	sc := echo.StartConfig{
		Address:         config.GetServer(),
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 10 * time.Second,
		BeforeServeFunc: func(s *http.Server) error {
			s.ReadHeaderTimeout = 10 * time.Second
			s.ReadTimeout = 30 * time.Second
			s.IdleTimeout = 120 * time.Second
			return nil
		},
	}
	if err := sc.Start(ctx, e); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
