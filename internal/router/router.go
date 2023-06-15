package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.unjx.de/flohoss/gobackup/internal/controller"
)

func InitRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	e.Validator = &CustomValidator{Validator: newValidator()}
	e.Renderer = initTemplates()

	return e
}

func SetupRoutes(e *echo.Echo, ctrl *controller.Controller) {
	static := e.Group("/static", longCacheLifetime)
	static.Static("/", "web/static")

	api := e.Group("/api/v1")
	api.GET("/logs", ctrl.GetLogs)
	api.POST("/job", ctrl.StartBackup)
	api.POST("/docker", ctrl.DockerRequest)
	api.POST("/restic", ctrl.ResticRequest)

	e.GET("/system", ctrl.RenderSystem)
	e.GET("/logs", ctrl.RenderLogs)
	e.GET("/tools", ctrl.RenderTools)
	e.GET("/sse", echo.WrapHandler(http.HandlerFunc(ctrl.SSE.ServeHTTP)))

	jobs := e.Group("/jobs")
	jobs.GET("", ctrl.RenderJobs)
	jobs.POST("", ctrl.CreateJobConfiguration)
	jobs.DELETE("/:id", ctrl.DeleteJobConfiguration)
	jobs.GET("/form", ctrl.RenderJobForm)

	remotes := e.Group("/remotes")
	remotes.GET("", ctrl.RenderRemotes)
	remotes.POST("", ctrl.CreateRemoteConfiguration)
	remotes.DELETE("/:id", ctrl.DeleteRemoteConfiguration)
	remotes.GET("/form", ctrl.RenderRemoteForm)

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})
	e.RouteNotFound("*", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/jobs")
	})
}
