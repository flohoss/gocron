package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupRoutes(e *echo.Echo, jh *JobHandler) {
	e.GET("/", jh.listHandler)
	e.GET("/:name", jh.jobHandler)

	api := e.Group("/api")
	api.GET("/events", jh.JobService.GetHandler())

	jobs := api.Group("/jobs")
	jobs.GET("", jh.listHandler)
	jobs.GET("/:name", jh.jobHandler)
	jobs.POST("", jh.executeJobsHandler)
	jobs.POST("/:name", jh.executeJobHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
