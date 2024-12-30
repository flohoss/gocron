package main

import (
	"context"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gitlab.unjx.de/flohoss/gobackup/services"
)

type config struct {
	DB_NAME string `env:"DB_NAME" envDefault:"db.sqlite"`
}

const (
	configFolder = "./config/"
)

func init() {
	os.Mkdir(configFolder, os.ModePerm)
}

func main() {
	var cfg config
	cfg, err := env.ParseAs[config]()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Debug = false

	e.Static("/static", "assets")
	e.Use(middleware.Logger())

	db, err := services.NewJobService(configFolder + cfg.DB_NAME)
	if err != nil {
		e.Logger.Error(err)
	}

	jobs, err := db.ListJobs(context.Background())
	if err != nil {
		e.Logger.Error(err)
	}

	log.Info("jobs", jobs)

	e.Logger.Fatal(e.Start(":8156"))
}
