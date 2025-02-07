package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "gitlab.unjx.de/flohoss/gobackup/docs"

	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/handlers"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"gitlab.unjx.de/flohoss/gobackup/internal/notify"
	"gitlab.unjx.de/flohoss/gobackup/internal/scheduler"
	"gitlab.unjx.de/flohoss/gobackup/services"
)

const (
	configFolder = "./config/"
)

func init() {
	os.Mkdir(configFolder, os.ModePerm)
}

// @title			No Backup No Mercy API
// @version		1.0
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath		/api
func main() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = false

	env, err := env.Parse()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Static("/static", "assets")
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "events") || strings.Contains(c.Path(), "docs")
		},
	}))

	cfg, err := config.New(configFolder + "config.yml")
	if err != nil {
		e.Logger.Fatal(err)
	}

	c := scheduler.New()
	n := notify.New(env.NtfyUrl, env.NtfyTopic, env.NtfyToken)

	js, err := services.NewJobService(configFolder+"db.sqlite", cfg, c, n)
	if err != nil {
		e.Logger.Fatal(err)
	}
	jh := handlers.NewJobHandler(js, cfg)

	handlers.SetupRoutes(e, jh)

	e.Logger.Fatal(e.Start(":8156"))
}
