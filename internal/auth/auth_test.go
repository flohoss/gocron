package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc/oidctest"
	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
)

func newTestAuth(t *testing.T) (*Auth, *jobs.Queries) {
	t.Helper()

	service, queries, _ := newTestAuthWithDB(t)

	return service, queries
}

func newTestAuthWithDB(t *testing.T) (*Auth, *jobs.Queries, *sql.DB) {
	t.Helper()

	loadConfig(t, true)

	db := newTestDB(t)
	queries := jobs.New(db)

	return &Auth{sessions: queries}, queries, db
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "auth.sqlite"))
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		jti TEXT PRIMARY KEY,
		subject TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		expires_at INTEGER NOT NULL
	)`); err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	return db
}

func persistSession(t *testing.T, queries *jobs.Queries, jti string, expiresAt time.Time) {
	t.Helper()

	if err := queries.CreateSession(context.Background(), jobs.CreateSessionParams{
		Jti:       jti,
		Subject:   "user-subject",
		CreatedAt: time.Now().UnixMilli(),
		ExpiresAt: expiresAt.UnixMilli(),
	}); err != nil {
		t.Fatalf("failed to persist session: %v", err)
	}
}

func serveWithSession(service *Auth, target, cookie string) *httptest.ResponseRecorder {
	router := echo.New()
	router.Use(service.Middleware())
	router.GET(target, func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, target, nil)
	if cookie != "" {
		req.Header.Set(echo.HeaderCookie, cookie)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestNew_DisabledWhenConfigDisabled(t *testing.T) {
	loadConfig(t, false)

	service := New(nil)

	if service.Enabled() {
		t.Fatal("expected auth to be disabled")
	}
}

// Enabled() reads the live config so toggling auth.oidc.enabled in the file
// takes effect without a restart.
func TestEnabled_FollowsConfigReload(t *testing.T) {
	loadConfig(t, true)

	service := New(nil)
	defer service.Shutdown()

	if !service.Enabled() {
		t.Fatal("expected auth to be enabled")
	}

	loadConfig(t, false)

	if service.Enabled() {
		t.Fatal("expected auth to be disabled after the reload")
	}
}

// Discovery is warmed in the background so the first login does not pay for the
// round trip, without blocking startup when the provider is slow or down.
func TestNew_WarmsProviderInBackground(t *testing.T) {
	oidcServer := &oidctest.Server{}
	server := httptest.NewServer(oidcServer)
	defer server.Close()
	oidcServer.SetIssuer(server.URL)

	loadConfigWithIssuer(t, server.URL)

	service := New(nil)
	defer service.Shutdown()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		service.providerMu.Lock()
		resolved := service.provider != nil
		service.providerMu.Unlock()

		if resolved {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("expected the provider to be discovered in the background")
}

// A provider that cannot be discovered degrades into a retryable error instead
// of a cached failure, so a provider that comes back later needs no restart.
func TestResolve_RetriesAfterAFailedDiscovery(t *testing.T) {
	loadConfig(t, true)

	service := &Auth{}

	if _, err := service.resolve(); !errors.Is(err, ErrProviderUnavailable) {
		t.Fatalf("expected ErrProviderUnavailable, got %v", err)
	}

	oidcServer := &oidctest.Server{}
	server := httptest.NewServer(oidcServer)
	defer server.Close()
	oidcServer.SetIssuer(server.URL)

	loadConfigWithIssuer(t, server.URL)

	resolved, err := service.resolve()
	if err != nil {
		t.Fatalf("expected discovery to succeed after the provider came back: %v", err)
	}
	if resolved == nil || resolved.verifier == nil {
		t.Fatal("expected a resolved provider")
	}
}

// The second attempt has to hit the network again rather than replay the first
// failure. The issuer is unchanged here, which the old cached error blocked.
func TestResolve_RetriesTheSameIssuerAfterAFailure(t *testing.T) {
	var healthy atomic.Bool

	oidcServer := &oidctest.Server{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !healthy.Load() {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		oidcServer.ServeHTTP(w, r)
	}))
	defer server.Close()
	oidcServer.SetIssuer(server.URL)

	loadConfigWithIssuer(t, server.URL)

	service := &Auth{}

	if _, err := service.resolve(); !errors.Is(err, ErrProviderUnavailable) {
		t.Fatalf("expected ErrProviderUnavailable, got %v", err)
	}

	healthy.Store(true)

	resolved, err := service.resolve()
	if err != nil {
		t.Fatalf("expected the same issuer to be retried: %v", err)
	}
	if resolved == nil {
		t.Fatal("expected a resolved provider")
	}
}

func TestNew_RequiresProviderFieldsWhenEnabled(t *testing.T) {
	v := baseConfig()
	v.Set("auth.oidc.enabled", true)

	if err := config.ValidateAndLoadConfig(v); err == nil {
		t.Fatal("expected validation error for missing OIDC fields, got nil")
	}
}

func TestMiddleware_DisabledLetsEverythingThrough(t *testing.T) {
	loadConfig(t, false)

	service := &Auth{}

	rec := serveWithSession(service, "/api/jobs", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected open access when disabled, got %d", rec.Code)
	}
}

func TestMiddleware_RejectsMissingCookieWithUnauthorized(t *testing.T) {
	service, _ := newTestAuth(t)

	rec := serveWithSession(service, "/api/jobs", "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without session, got %d", rec.Code)
	}
}

func TestMiddleware_AcceptsValidSession(t *testing.T) {
	service, queries := newTestAuth(t)
	expiresAt := time.Now().Add(time.Hour)
	token, jti := randomTokenForTest(t), randomTokenForTest(t)
	_ = token
	persistSession(t, queries, jti, expiresAt)

	rec := serveWithSession(service, "/api/jobs", SessionCookieName+"="+jti)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid session, got %d", rec.Code)
	}
}

// Deleting the session is the reason sessions live in the database: a cookie
// value that is still present must stop working once logged out.
func TestMiddleware_RejectsDeletedSession(t *testing.T) {
	service, queries := newTestAuth(t)
	expiresAt := time.Now().Add(time.Hour)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, expiresAt)

	if err := queries.DeleteSession(context.Background(), jti); err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	rec := serveWithSession(service, "/api/jobs", SessionCookieName+"="+jti)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for deleted session, got %d", rec.Code)
	}
}

func TestMiddleware_RejectsExpiredSession(t *testing.T) {
	service, queries := newTestAuth(t)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, time.Now().Add(-time.Minute))

	rec := serveWithSession(service, "/api/jobs", SessionCookieName+"="+jti)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired session, got %d", rec.Code)
	}
}

func TestAuthenticate_ReturnsPersistedIdentity(t *testing.T) {
	service, queries := newTestAuth(t)
	expiresAt := time.Now().Add(time.Hour)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, expiresAt)

	router := echo.New()
	var user *User
	router.GET("/api/jobs", func(c *echo.Context) error {
		var err error
		user, err = service.Authenticate(c)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	req.Header.Set(echo.HeaderCookie, SessionCookieName+"="+jti)
	router.ServeHTTP(httptest.NewRecorder(), req)

	if user == nil || user.Subject != "user-subject" {
		t.Fatalf("unexpected authenticated user: %+v", user)
	}
}

func TestLogout_DeletesPersistedSession(t *testing.T) {
	service, queries, db := newTestAuthWithDB(t)
	expiresAt := time.Now().Add(time.Hour)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, expiresAt)

	router := echo.New()
	router.POST("/api/auth/logout", func(c *echo.Context) error {
		if err := service.Logout(c); err != nil {
			t.Fatalf("unexpected logout error: %v", err)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.Header.Set(echo.HeaderCookie, SessionCookieName+"="+jti)
	router.ServeHTTP(httptest.NewRecorder(), req)

	var remaining int
	if err := db.QueryRow("SELECT COUNT(*) FROM sessions WHERE jti = ?", jti).Scan(&remaining); err != nil {
		t.Fatalf("failed to count sessions: %v", err)
	}
	if remaining != 0 {
		t.Fatalf("expected the session row to be deleted, found %d row(s)", remaining)
	}
}

func TestDeleteExpiredSessions_RemovesOnlyExpiredRows(t *testing.T) {
	service, queries := newTestAuth(t)

	expired := time.Now().Add(-time.Hour)
	jtiExpired := randomTokenForTest(t)
	persistSession(t, queries, jtiExpired, expired)

	valid := time.Now().Add(time.Hour)
	jtiValid := randomTokenForTest(t)
	persistSession(t, queries, jtiValid, valid)

	service.deleteExpiredSessions()

	if _, err := queries.GetActiveSession(context.Background(), jobs.GetActiveSessionParams{
		Jti:       jtiExpired,
		ExpiresAt: time.Now().UnixMilli(),
	}); err == nil {
		t.Fatal("expected expired session to be deleted")
	}
	if _, err := queries.GetActiveSession(context.Background(), jobs.GetActiveSessionParams{
		Jti:       jtiValid,
		ExpiresAt: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatalf("expected valid session to survive cleanup: %v", err)
	}
}

func TestStartSession_PersistsAndSetsCookie(t *testing.T) {
	service, queries := newTestAuth(t)

	router := echo.New()
	router.GET("/api/auth/callback", func(c *echo.Context) error {
		return service.StartSession(c, "user-subject", time.Now().Add(time.Hour))
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/auth/callback", nil))

	sessionCookie := rec.Result().Cookies()[0]
	if sessionCookie.Name != SessionCookieName || sessionCookie.Value == "" {
		t.Fatalf("expected a session cookie with a token, got %+v", sessionCookie)
	}

	if _, err := queries.GetActiveSession(context.Background(), jobs.GetActiveSessionParams{
		Jti:       sessionCookie.Value,
		ExpiresAt: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatalf("expected the session cookie value to be a persisted session: %v", err)
	}
}

// Rotating cookie flags or the session TTL must not require a restart: these
// helpers read the reloaded config at call time.
func TestSessionCookie_FollowsConfigReload(t *testing.T) {
	service := &Auth{}

	loadConfig(t, true)
	if got := service.SessionTTL(); got != 24*time.Hour {
		t.Fatalf("expected the hydrated default TTL, got %v", got)
	}

	expiresAt := time.Now().Add(time.Hour)
	cookie := service.SessionCookie("token", expiresAt)

	if cookie.Name != SessionCookieName {
		t.Fatalf("unexpected cookie name: %q", cookie.Name)
	}
	if !cookie.HttpOnly {
		t.Fatal("expected session cookie to be HttpOnly")
	}
	if cookie.Secure {
		t.Fatal("expected an insecure cookie before enabling cookie_secure")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite Lax, got %v", cookie.SameSite)
	}
	if !cookie.Expires.Equal(expiresAt) {
		t.Fatalf("expected cookie expiry %v, got %v", expiresAt, cookie.Expires)
	}

	v := baseConfig()
	v.Set("auth.oidc.enabled", true)
	v.Set("auth.oidc.issuer_url", "https://sso.example.com")
	v.Set("auth.oidc.client_id", "gocron")
	v.Set("auth.oidc.client_secret", "secret")
	v.Set("auth.oidc.cookie_secure", true)
	v.Set("auth.oidc.session_ttl", 2*time.Hour)
	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to reload test config: %v", err)
	}

	if got := service.SessionTTL(); got != 2*time.Hour {
		t.Fatalf("expected the reloaded TTL, got %v", got)
	}
	if cookie := service.SessionCookie("token", expiresAt); !cookie.Secure {
		t.Fatal("expected the session cookie to honour the reloaded cookie_secure")
	}
}

func TestClearSessionCookie_ExpiresImmediately(t *testing.T) {
	service := &Auth{}

	if cookie := service.ClearSessionCookie(); cookie.MaxAge != -1 {
		t.Fatalf("expected MaxAge -1 to delete the cookie, got %d", cookie.MaxAge)
	}
}

func TestRandomToken_ReturnsDistinctHexValues(t *testing.T) {
	first, err := randomToken(sessionTokenSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := randomToken(sessionTokenSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first == "" || second == "" {
		t.Fatal("expected non-empty tokens")
	}
	if first == second {
		t.Fatal("expected each token to be unique")
	}
	if len(first) != 2*sessionTokenSize {
		t.Fatalf("expected %d hex characters, got %d", 2*sessionTokenSize, len(first))
	}
}

func randomTokenForTest(t *testing.T) string {
	t.Helper()

	token, err := randomToken(sessionTokenSize)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return token
}

// The provider compares redirect_uri against the URL registered for the client,
// so it has to be derived from the request the browser actually made.
func TestCallbackURL_DerivedFromRequest(t *testing.T) {
	cases := []struct {
		name   string
		target string
		proto  string
		want   string
	}{
		{name: "direct http", target: "http://gocron.example.com:8156/api/auth/login", want: "http://gocron.example.com:8156/api/auth/callback"},
		{name: "TLS", target: "https://gocron.example.com/api/auth/login", want: "https://gocron.example.com/api/auth/callback"},
		{name: "forwarded proto", target: "http://gocron:8156/api/auth/login", proto: "https", want: "https://gocron:8156/api/auth/callback"},
	}

	for _, tc := range cases {
		router := echo.New()
		var got string
		router.GET("/api/auth/login", func(c *echo.Context) error {
			got = callbackURL(c)
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, tc.target, nil)
		if tc.proto != "" {
			req.Header.Set(echo.HeaderXForwardedProto, tc.proto)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		if got != tc.want {
			t.Errorf("%s: callbackURL() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func loadConfig(t *testing.T, enabled bool) {
	t.Helper()

	v := baseConfig()
	if enabled {
		v.Set("auth.oidc.enabled", true)
		v.Set("auth.oidc.issuer_url", "https://sso.example.com")
		v.Set("auth.oidc.client_id", "gocron")
		v.Set("auth.oidc.client_secret", "secret")
	}

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
}

func loadConfigWithIssuer(t *testing.T, issuer string) {
	t.Helper()

	v := baseConfig()
	v.Set("auth.oidc.enabled", true)
	v.Set("auth.oidc.issuer_url", issuer)
	v.Set("auth.oidc.client_id", "gocron")
	v.Set("auth.oidc.client_secret", "secret")

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
}

func baseConfig() *viper.Viper {
	v := viper.New()
	v.Set("time_zone", "UTC")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("jobs", []map[string]any{{
		"name":     "Auth Test Job",
		"commands": []string{"echo test"},
	}})
	return v
}
