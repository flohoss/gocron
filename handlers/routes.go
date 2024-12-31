package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, jh *JobHandler) {
	e.GET("/", jh.listHandler)
	e.GET("/:name", jh.jobHandler)
}
