package handlers

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/flohoss/gocron/internal/webui"
	"github.com/labstack/echo/v4"
)

func parseIndexTemplate() *template.Template {
	if distFS, ok := webui.DistFS(); ok {
		return template.Must(template.ParseFS(distFS, "index.html"))
	}

	return template.Must(template.ParseGlob("web/index.html"))
}

func registerStaticRoutes(e *echo.Echo) {
	if distFS, ok := webui.DistFS(); ok {
		registerEmbeddedStatic(e, distFS, "/assets", "assets")
		registerEmbeddedStatic(e, distFS, "/static", "static")
		return
	}

	assets := e.Group("/assets", longCacheLifetime)
	assets.Static("/", "web/assets")

	static := e.Group("/static", longCacheLifetime)
	static.Static("/", "web/static")
}

func registerEmbeddedStatic(e *echo.Echo, distFS fs.FS, routePrefix string, dir string) {
	subFS, err := fs.Sub(distFS, dir)
	if err != nil {
		panic(err)
	}

	handler := http.StripPrefix(routePrefix, http.FileServerFS(subFS))
	group := e.Group(routePrefix, longCacheLifetime)
	group.GET("/*", echo.WrapHandler(handler))
}
