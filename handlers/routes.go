package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/internal/events"
)

func SetupRoutes(e *echo.Echo, jh *JobHandler) {
	e.GET("/", jh.listHandler)
	e.GET("/:name", jh.jobHandler)

	api := e.Group("/api")
	api.GET("/events", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jh.JobService.GetEvents().SendGlobal(&events.GlobalInfo{
			Idle: jh.JobService.IsIdle(),
		})
		jh.JobService.GetEvents().SSE.ServeHTTP(w, r)
	})))

	jobs := api.Group("/jobs")
	jobs.POST("", jh.executeJobsHandler)
	jobs.POST("/:name", jh.executeJobHandler)

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
