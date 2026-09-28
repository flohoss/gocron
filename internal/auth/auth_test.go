package auth

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
)

const testSigningKey = "test-signing-key"

func newTestAuth(t *testing.T) (*Auth, *jobs.Queries) {
	t.Helper()

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "auth.sqlite"))
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		jti TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		expires_at INTEGER NOT NULL,
		revoked INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	service := &Auth{
		enabled:  true,
		sessions: jobs.New(db),
	}

	v := baseConfig()
	v.Set("auth.oidc.session_ttl", time.Hour)
	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	return service, jobs.New(db)
}

func persistSession(t *testing.T, queries *jobs.Queries, jti string, expiresAt time.Time) {
	t.Helper()

	if err := queries.CreateSession(context.Background(), jobs.CreateSessionParams{
		Jti:       jti,
		Email:     "user@example.com",
		Username:  "user",
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

	service, err := New(nil)
	if err != nil {
		t.Fatalf("expected no error when SSO disabled, got: %v", err)
	}
	if service.Enabled() {
		t.Fatal("expected auth to be disabled")
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

// Revocation is the reason sessions live in the database: a cookie value that is
// still present must stop working once logged out.
func TestMiddleware_RejectsRevokedSession(t *testing.T) {
	service, queries := newTestAuth(t)
	expiresAt := time.Now().Add(time.Hour)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, expiresAt)

	if err := queries.RevokeSession(context.Background(), jti); err != nil {
		t.Fatalf("failed to revoke session: %v", err)
	}

	rec := serveWithSession(service, "/api/jobs", SessionCookieName+"="+jti)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for revoked session, got %d", rec.Code)
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

	if user == nil || user.Email != "user@example.com" || user.Username != "user" {
		t.Fatalf("unexpected authenticated user: %+v", user)
	}
}

func TestLogout_RevokesPersistedSession(t *testing.T) {
	service, queries := newTestAuth(t)
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

	if _, err := queries.GetActiveSession(context.Background(), jobs.GetActiveSessionParams{
		Jti:       jti,
		ExpiresAt: time.Now().UnixMilli(),
	}); err == nil {
		t.Fatal("expected revoked session to be rejected, got an active session")
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
		return service.StartSession(c, &User{Email: "user@example.com", Username: "user"}, time.Now().Add(time.Hour))
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
	v.Set("auth.oidc.auth_url", "https://sso.example.com/authorize")
	v.Set("auth.oidc.token_url", "https://sso.example.com/token")
	v.Set("auth.oidc.jwks_url", "https://sso.example.com/keys")
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
		v.Set("auth.oidc.auth_url", "https://sso.example.com/authorize")
		v.Set("auth.oidc.token_url", "https://sso.example.com/token")
		v.Set("auth.oidc.jwks_url", "https://sso.example.com/keys")
		v.Set("auth.oidc.client_id", "gocron")
		v.Set("auth.oidc.client_secret", "secret")
	}

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

// Both the spec's boolean and the string form must be accepted; anything
// ambiguous must not count as a verified email.
func TestIsEmailVerified(t *testing.T) {
	cases := []struct {
		name  string
		claim any
		want  bool
	}{
		{name: "bool true", claim: true, want: true},
		{name: "bool false", claim: false, want: false},
		{name: "string true", claim: "true", want: true},
		{name: "string false", claim: "false", want: false},
		{name: "uppercase string true", claim: "TRUE", want: true},
		{name: "numeric one", claim: 1, want: false},
		{name: "missing", claim: nil, want: false},
		{name: "unparsable string", claim: "yes", want: false},
		{name: "empty string", claim: "", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isEmailVerified(tc.claim); got != tc.want {
				t.Fatalf("isEmailVerified(%#v) = %v, want %v", tc.claim, got, tc.want)
			}
		})
	}
}
