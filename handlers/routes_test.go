package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flohoss/gocron/config"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
)

func TestBuildIPExtractor_DefaultsToDirectConnection(t *testing.T) {
	extractor := buildIPExtractor(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.7:1234"
	req.Header.Set(echo.HeaderXForwardedFor, "203.0.113.9")

	// Without trusted proxies a spoofed X-Forwarded-For must be ignored.
	if got := extractor(req); got != "198.51.100.7" {
		t.Fatalf("expected direct IP, got %q", got)
	}
}

func TestBuildIPExtractor_UsesForwardedForFromTrustedProxy(t *testing.T) {
	extractor := buildIPExtractor([]string{"10.0.0.0/8"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.5:1234"
	req.Header.Set(echo.HeaderXForwardedFor, "203.0.113.9")

	if got := extractor(req); got != "203.0.113.9" {
		t.Fatalf("expected forwarded client IP, got %q", got)
	}
}

func TestBuildIPExtractor_IgnoresForwardedForFromUntrustedPeer(t *testing.T) {
	extractor := buildIPExtractor([]string{"10.0.0.0/8"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.7:1234"
	req.Header.Set(echo.HeaderXForwardedFor, "203.0.113.9")

	if got := extractor(req); got != "198.51.100.7" {
		t.Fatalf("expected direct IP for untrusted peer, got %q", got)
	}
}

func TestBuildIPExtractor_IgnoresInvalidEntries(t *testing.T) {
	extractor := buildIPExtractor([]string{"not-a-cidr", "10.0.0.0/8"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.5:1234"
	req.Header.Set(echo.HeaderXForwardedFor, "203.0.113.9")

	if got := extractor(req); got != "203.0.113.9" {
		t.Fatalf("expected valid entry to still apply, got %q", got)
	}
}

func TestBuildRateLimitMiddleware_DisabledIsNoop(t *testing.T) {
	limiter := buildRateLimitMiddleware()

	handler := limiter(func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	router := echo.New()
	router.GET("/", handler)

	for i := 0; i < 50; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "198.51.100.7:1234"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200 when limiter disabled, got %d", i, rec.Code)
		}
	}
}

func TestBuildRateLimitMiddleware_BlocksAboveBurst(t *testing.T) {
	loadTestServerConfig(t, true, 1, 2)

	limiter := buildRateLimitMiddleware()
	handler := limiter(func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	router := echo.New()
	router.GET("/", handler)

	var blocked int
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "198.51.100.7:1234"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code == http.StatusTooManyRequests {
			blocked++
		}
	}

	if blocked == 0 {
		t.Fatal("expected the rate limiter to reject requests above the burst")
	}
}

// Allowed origins are re-read per request so restricting them only needs the
// config file reload, not a restart.
func TestBuildCORSMiddleware_FollowsConfigReload(t *testing.T) {
	loadTestServerConfig(t, false, 0, 0)

	router := echo.New()
	router.Use(buildCORSMiddleware())
	router.GET("/", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	serveOrigin := func(origin string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderOrigin, origin)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}

	if rec := serveOrigin("https://evil.example.com"); rec.Header().Get(echo.HeaderAccessControlAllowOrigin) != "*" {
		t.Fatal("expected wildcard allow-origin before the reload")
	}

	v := viper.New()
	v.Set("time_zone", "UTC")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("server.cors.allow_origins", []string{"https://app.example.com"})
	v.Set("jobs", []map[string]any{{
		"name":     "Router Test Job",
		"commands": []string{"echo test"},
	}})
	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to reload test config: %v", err)
	}

	if got := serveOrigin("https://app.example.com").Header().Get(echo.HeaderAccessControlAllowOrigin); got != "https://app.example.com" {
		t.Fatalf("expected the reloaded origin to be allowed, got %q", got)
	}
	if got := serveOrigin("https://evil.example.com").Header().Get(echo.HeaderAccessControlAllowOrigin); got != "" {
		t.Fatalf("expected the unknown origin to be rejected after reload, got %q", got)
	}
}

// Preflight advertises Echo's default method list rather than the methods this
// API registers. The SPA catch-all means unmatched requests resolve to a
// not-found handler, so Echo never records the route's Allow value in the
// context and CORS falls back to its defaults. DELETE and PATCH therefore appear
// in the header even though no route implements them.
func TestBuildCORSMiddleware_AdvertisesEchoDefaultMethods(t *testing.T) {
	loadTestServerConfig(t, false, 0, 0)

	router := echo.New()
	router.Use(buildCORSMiddleware())
	router.GET("/api/jobs", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	router.RouteNotFound("*", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/jobs", nil)
	req.Header.Set(echo.HeaderOrigin, "https://example.com")
	req.Header.Set(echo.HeaderAccessControlRequestMethod, http.MethodGet)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	got := rec.Header().Get(echo.HeaderAccessControlAllowMethods)
	want := "GET,HEAD,PUT,PATCH,POST,DELETE"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func loadTestServerConfig(t *testing.T, rateLimitEnabled bool, rate float64, burst int) {
	t.Helper()

	v := viper.New()
	v.Set("time_zone", "UTC")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("server.rate_limit.enabled", rateLimitEnabled)
	v.Set("server.rate_limit.rate", rate)
	v.Set("server.rate_limit.burst", burst)
	v.Set("jobs", []map[string]any{{
		"name":     "Router Test Job",
		"commands": []string{"echo test"},
	}})

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
}
