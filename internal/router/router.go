package router

import (
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	e.Use(middleware.Gzip())

	e.Validator = &CustomValidator{Validator: newValidator()}

	return e
}

func SetupRoutes(e *echo.Echo, ctrl *controller.Controller, env *env.Config) {
	api := e.Group("/api")

	jobs := api.Group("/jobs")
	jobs.GET("", ctrl.GetJobs)
	jobs.GET("/:id", ctrl.GetJob)
	jobs.PUT("", ctrl.UpdateJob)
	jobs.POST("", ctrl.CreateJob)
	jobs.DELETE("/:id", ctrl.DeleteJob)

	retentionPolicy := api.Group("/retention_policies")
	retentionPolicy.GET("", ctrl.GetRetentionPolicies)

	compressionTypes := api.Group("/compression_types")
	compressionTypes.GET("", ctrl.GetCompressionTypes)

	system := api.Group("/system")
	system.GET("", ctrl.GetSystem)

	if env.SwaggerHost != "" {
		docs.SwaggerInfo.Title = "GoBackup"
		docs.SwaggerInfo.Version = env.Version
		docs.SwaggerInfo.BasePath = "/api"
		parsed, _ := url.Parse(env.SwaggerHost)
		docs.SwaggerInfo.Host = parsed.Host

		api.GET("/swagger/*", echoSwagger.WrapHandler)
		zap.L().Info("swagger running", zap.String("url", env.SwaggerHost+docs.SwaggerInfo.BasePath+"/swagger/index.html"))
	}
}
