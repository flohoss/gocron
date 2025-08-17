package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
)

func longCacheLifetime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		return next(c)
	}
}

func SetupRouter(e *echo.Echo, jh *JobHandler, ch *CommandHandler) {
	config := huma.DefaultConfig("GoCron API", "1.0.0")
	config.OpenAPIPath = "/api/openapi"
	config.SchemasPath = "/api/schemas"
	h := humaecho.New(e, config)

	e.GET("/api/docs", func(ctx echo.Context) error {
		return ctx.HTML(http.StatusOK, `<!doctype html>
			<html>
				<head>
					<title>API Reference</title>
					<meta charset="utf-8" />
					<meta name="viewport" content="width=device-width, initial-scale=1" />
				</head>
				<body>
					<script id="api-reference" data-url="/api/openapi.json"></script>
					<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
				</body>
			</html>`,
		)
	})
	e.Renderer = initTemplates()

	e.GET("/api/events", jh.JobService.GetHandler())
	huma.Register(h, ch.executeCommandOperation(), ch.executeCommandHandler)
	huma.Register(h, jh.listJobsOperation(), jh.listJobsHandler)
	huma.Register(h, jh.listRunsOperation(), jh.listRunsHandler)
	huma.Register(h, jh.executeJobsOperation(), jh.executeJobsHandler)
	huma.Register(h, jh.executeJobOperation(), jh.executeJobHandler)
	huma.Register(h, jh.changeJobOperation(), jh.changeJobHandler)

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
