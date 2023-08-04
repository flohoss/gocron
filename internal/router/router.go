package router

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/docs"
	"gitlab.unjx.de/flohoss/gobackup/internal/controller"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"go.uber.org/zap"
)

func InitRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "sse") || strings.Contains(c.Path(), "swagger")
		},
	}))

	e.Validator = &CustomValidator{Validator: newValidator()}
	e.Renderer = initTemplates()

	return e
}

func SetupRoutes(e *echo.Echo, ctrl *controller.Controller, env *env.Config) {
	favicon := e.Group("/favicon", longCacheLifetime)
	favicon.Static("/", "web/favicon")
	e.Static("/assets", "web/assets")

	api := e.Group("/api")
	api.GET("/sse", echo.WrapHandler(http.HandlerFunc(database.SSE.ServeHTTP)))

	jobs := api.Group("/jobs")
	jobs.GET("", ctrl.GetJobs)
	jobs.GET("/:id", ctrl.GetJob)
	jobs.PUT("", ctrl.UpdateJob)
	jobs.POST("", ctrl.CreateJob)
	jobs.DELETE("/:id", ctrl.DeleteJob)
	jobs.GET("/:id/runs", ctrl.GetRuns)

	commands := api.Group("/commands")
	commands.POST("", ctrl.RunCommand)

	system := api.Group("/system")
	system.GET("", ctrl.GetSystem)
	system.GET("/stats", ctrl.GetSystemStats)
	system.GET("/logs", ctrl.GetSystemLogs)

	if env.SwaggerHost != "" {
		docs.SwaggerInfo.Title = "GoBackup"
		docs.SwaggerInfo.Version = env.Version
		docs.SwaggerInfo.BasePath = "/api"
		parsed, _ := url.Parse(env.SwaggerHost)
		docs.SwaggerInfo.Host = parsed.Host

		api.GET("/swagger/*", echoSwagger.WrapHandler)
		zap.L().Info("swagger running", zap.String("url", env.SwaggerHost+docs.SwaggerInfo.BasePath+"/swagger/index.html"))
	}

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})
	e.RouteNotFound("*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
}
