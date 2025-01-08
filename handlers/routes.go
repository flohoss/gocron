package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, jh *JobHandler) {
	e.GET("/", jh.listHandler)
	e.GET("/:name", jh.jobHandler)
	e.GET("/:name", jh.jobHandler)

	api := e.Group("/api")

	jobs := api.Group("/jobs")
	jobs.POST("", jh.executeJobsHandler)
}
