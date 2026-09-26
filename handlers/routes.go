package handlers

import (
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/internal/buildinfo"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func longCacheLifetime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		return next(c)
	}
}

func healthHandler(c *echo.Context) error {
	return c.String(http.StatusOK, ".")
}

func InitRouter() *echo.Echo {
	e := echo.NewWithConfig(echo.Config{
		Logger:      slog.Default(),
		IPExtractor: buildIPExtractor(config.GetTrustedProxies()),
	})

	e.Use(middleware.Recover())
	e.Use(buildCORSMiddleware())
	e.Use(buildRateLimitMiddleware())
	e.Use(buildRequestLoggerMiddleware())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c *echo.Context) bool {
			return strings.Contains(c.Path(), "events")
		},
	}))

	e.Renderer = initTemplates()

	return e
}

func buildRequestLoggerMiddleware() echo.MiddlewareFunc {
	if config.GetLogLevel() != slog.LevelDebug {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return middleware.RequestLogger()
}

func buildIPExtractor(trustedProxies []string) echo.IPExtractor {
	if len(trustedProxies) == 0 {
		return echo.ExtractIPDirect()
	}

	options := []echo.TrustOption{
		echo.TrustLoopback(false),
		echo.TrustLinkLocal(false),
		echo.TrustPrivateNet(false),
	}
	for _, proxy := range trustedProxies {
		_, network, err := net.ParseCIDR(proxy)
		if err != nil {
			slog.Error("Ignoring invalid trusted proxy entry", "entry", proxy, "error", err)
			continue
		}
		options = append(options, echo.TrustIPRange(network))
	}

	return echo.ExtractIPFromXFFHeader(options...)
}

func buildCORSMiddleware() echo.MiddlewareFunc {
	settings := config.GetCORSSettings()

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: settings.AllowOrigins,
	})
}

func buildRateLimitMiddleware() echo.MiddlewareFunc {
	settings := config.GetRateLimitSettings()
	if !settings.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:  settings.Rate,
		Burst: settings.Burst,
	})

	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store,
	})
}

func SetupRouter(e *echo.Echo, jh *JobHandler, ch *CommandHandler) {
	e.GET("/health", healthHandler)
	e.HEAD("/health", healthHandler)

	h := huma.DefaultConfig("GoCron API", buildinfo.Version)
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

	e.GET("/robots.txt", func(ctx *echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	registerStaticRoutes(e)

	e.RouteNotFound("*", func(ctx *echo.Context) error {
		return ctx.Render(http.StatusOK, "index.html", nil)
	})
}
