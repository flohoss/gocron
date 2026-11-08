package handlers

import (
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/internal/buildinfo"
	"github.com/flohoss/gosso/v3"
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
		// Group middleware must stay scoped to the routes registered through
		// the group. Without this flag Echo auto-registers 404 catch-alls for
		// every group with middleware, which would run the session check on
		// every unmatched request including the SPA shell.
		NoGroupAutoRegister404Routes: true,
	})

	e.Use(middleware.Recover())
	e.Use(gosso.EchoContextMiddleware)
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
	requestLogger := middleware.RequestLogger()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		withLogger := requestLogger(next)
		return func(c *echo.Context) error {
			if config.GetLogLevel() != slog.LevelDebug {
				return next(c)
			}

			return withLogger(c)
		}
	}
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
	return middleware.CORSWithConfig(middleware.CORSConfig{
		UnsafeAllowOriginFunc: func(c *echo.Context, origin string) (string, bool, error) {
			allowed := config.GetCORSSettings().AllowOrigins
			if slices.Contains(allowed, "*") {
				return "*", true, nil
			}
			if slices.Contains(allowed, origin) {
				return origin, true, nil
			}
			return "", false, nil
		},
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

func SetupRouter(e *echo.Echo, auth *gosso.Auth, jh *JobHandler, ch *CommandHandler) {
	e.GET("/health", healthHandler)
	e.HEAD("/health", healthHandler)

	h := huma.DefaultConfig("GoCron API", buildinfo.Version)
	h.OpenAPIPath = "/api/openapi"
	h.DocsPath = "/api/docs"
	h.SchemasPath = "/api/schemas"
	h.Servers = []*huma.Server{{URL: ""}}

	// Only the routes registered on the returned API carry the session check.
	// Echo applies group middleware by route, so the public spec, docs, schemas,
	// SPA shell, assets and the auth lifecycle stay open without any skipper.
	api := auth.Mount(e, h)

	e.GET("/api/events", auth.Middleware()(jh.JobService.GetHandler()))
	huma.Register(api, ch.executeCommandOperation(), ch.executeCommandHandler)
	huma.Register(api, jh.listJobsOperation(), jh.listJobsHandler)
	huma.Register(api, jh.listRunsOperation(), jh.listRunsHandler)
	huma.Register(api, jh.executeJobsOperation(), jh.executeJobsHandler)
	huma.Register(api, jh.executeJobOperation(), jh.executeJobHandler)
	huma.Register(api, jh.changeJobOperation(), jh.changeJobHandler)

	e.GET(gosso.LoginPagePath, func(c *echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/robots.txt", func(ctx *echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	registerStaticRoutes(e)
	registerFallbackRoutes(e, auth)
}

func registerFallbackRoutes(e *echo.Echo, auth *gosso.Auth) {
	e.RouteNotFound("/api/*", func(ctx *echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	})

	e.RouteNotFound("/*", func(ctx *echo.Context) error {
		if ctx.Request().Method != http.MethodGet {
			return echo.NewHTTPError(http.StatusNotFound, "Not found")
		}

		if auth.Enabled() {
			if _, err := auth.Authenticate(ctx); err != nil {
				return ctx.Redirect(http.StatusFound, gosso.LoginPagePath)
			}
		}

		return ctx.Render(http.StatusOK, "index.html", nil)
	})
}
