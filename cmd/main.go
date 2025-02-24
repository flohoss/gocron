package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

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

func main() {
	env, err := env.Parse()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	cfg, err := config.New(configFolder + "config.yml")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	c := scheduler.New()
	n := notify.New(env.NtfyUrl, env.NtfyTopic, env.NtfyToken)

	js, err := services.NewJobService(configFolder+"db.sqlite", cfg, c, n)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	jh := handlers.NewJobHandler(js)

	e := handlers.SetupRouter(jh)

	if err := e.Start(":8156"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
