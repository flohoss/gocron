package handlers

import (
	"net/http"
	"strings"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func longCacheLifetime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		return next(c)
	}
}

func SetupRouter(jh *JobHandler) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(echo.WrapMiddleware(chimiddleware.Heartbeat("/api/health")))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "events")
		},
	}))
	e.Renderer = initTemplates()

	api := e.Group("/api")
	api.GET("/events", jh.JobService.GetHandler())
	api.GET("/versions", jh.getVersions)

	jobs := api.Group("/jobs")
	jobs.GET("", jh.listHandler)
	jobs.GET("/:name", jh.jobHandler)
	jobs.POST("", jh.executeJobsHandler)
	jobs.POST("/:name", jh.executeJobHandler)

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	assets := e.Group("/assets", longCacheLifetime)
	assets.Static("/", "web/assets")

	favicon := e.Group("/static", longCacheLifetime)
	favicon.Static("/", "web/static")

	e.RouteNotFound("*", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "index.html", nil)
	})

	return e
}
