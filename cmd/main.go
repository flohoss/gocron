package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.unjx.de/flohoss/gobackup/config"
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

	_, err = services.NewJobService(configFolder+"db.sqlite", cfg)
	if err != nil {
		e.Logger.Error(err)
	}

	e.Logger.Fatal(e.Start(":8156"))
}
