package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func longCacheLifetime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		return next(c)
	}
}

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, ".")
}

func InitRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "events")
		},
	}))

	e.Renderer = initTemplates()

	return e
}

func SetupRouter(e *echo.Echo, jh *JobHandler, ch *CommandHandler) {
	e.GET("/health", healthHandler)
	e.HEAD("/health", healthHandler)

	h := huma.DefaultConfig("GoCron API", os.Getenv("APP_VERSION"))
	h.OpenAPIPath = "/api/openapi"
	h.DocsPath = "/api/docs"
	h.SchemasPath = "/api/schemas"
	humaAPI := humaecho.New(e, h)

	e.GET("/api/events", jh.JobService.GetHandler())
	huma.Register(humaAPI, ch.executeCommandOperation(), ch.executeCommandHandler)
	huma.Register(humaAPI, jh.listJobsOperation(), jh.listJobsHandler)
	huma.Register(humaAPI, jh.listRunsOperation(), jh.listRunsHandler)
	huma.Register(humaAPI, jh.executeJobsOperation(), jh.executeJobsHandler)
	huma.Register(humaAPI, jh.executeJobOperation(), jh.executeJobHandler)
	huma.Register(humaAPI, jh.changeJobOperation(), jh.changeJobHandler)

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
}
