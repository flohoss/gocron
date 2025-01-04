package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/handlers"
	"gitlab.unjx.de/flohoss/gobackup/internal/cron"
	"gitlab.unjx.de/flohoss/gobackup/services"
)

const (
	configFolder = "./config/"
)

func init() {
	os.Mkdir(configFolder, os.ModePerm)
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = false

	e.Static("/static", "assets")
	e.Use(middleware.Logger())

	cfg, err := config.New(configFolder + "config.yml")
	if err != nil {
		e.Logger.Error(err)
	}

	cron := cron.New()

	js, err := services.NewJobService(configFolder+"db.sqlite", cfg, cron)
	if err != nil {
		e.Logger.Error(err)
	}
	jh := handlers.NewJobHandler(js, cfg)

	handlers.SetupRoutes(e, jh)

	cron.Run()

	e.Logger.Fatal(e.Start(":8156"))
}
