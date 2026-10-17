package auth

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc/oidctest"
	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
	"github.com/flohoss/gosso/v2"
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

	return &Auth{sessions: queries, core: newCore(queries)}, queries, db
}

func newCore(queries *jobs.Queries) *gosso.Auth {
	return gosso.New(settings, NewDBSessions(queries), callbackPath, "")
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
		if service.core.Provider().Resolved() {
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

	service := New(nil)
	defer service.Shutdown()

	if _, err := service.core.Provider().Resolve(); !errors.Is(err, gosso.ErrProviderUnavailable) {
		t.Fatalf("expected ErrProviderUnavailable, got %v", err)
	}

	oidcServer := &oidctest.Server{}
	server := httptest.NewServer(oidcServer)
	defer server.Close()
	oidcServer.SetIssuer(server.URL)

	loadConfigWithIssuer(t, server.URL)

	if !service.core.Provider().Resolved() {
		if _, err := service.core.Provider().Resolve(); err != nil {
			t.Fatalf("expected discovery to succeed after the provider came back: %v", err)
		}
	}
}

// A server that accepts the connection but never responds would otherwise hang
// discovery forever. The timeout keeps the caller unblocked so login can be
// retried on the next request.
func TestResolve_TimesOutOnAStuckProvider(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open listener: %v", err)
	}
	t.Cleanup(func() { listener.Close() })

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				<-t.Context().Done()
				conn.Close()
			}()
		}
	}()

	loadConfigWithIssuer(t, "http://"+listener.Addr().String())

	service := New(nil)
	defer service.Shutdown()

	started := time.Now()
	_, err = service.core.Provider().Resolve()
	elapsed := time.Since(started)

	if !errors.Is(err, gosso.ErrProviderUnavailable) {
		t.Fatalf("expected ErrProviderUnavailable, got %v", err)
	}
	if elapsed > 15*time.Second {
		t.Fatalf("discovery took %v, expected it to give up instead of hanging", elapsed)
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

	service := New(nil)
	defer service.Shutdown()

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
	jti := randomTokenForTest(t)
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

func TestClearSession_DeletesPersistedSession(t *testing.T) {
	service, queries, db := newTestAuthWithDB(t)
	expiresAt := time.Now().Add(time.Hour)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, expiresAt)

	router := echo.New()
	router.POST("/api/auth/logout", func(c *echo.Context) error {
		service.ClearSession(c)
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

func TestClearSession_ExpiresCookie(t *testing.T) {
	service, queries, _ := newTestAuthWithDB(t)
	expiresAt := time.Now().Add(time.Hour)
	jti := randomTokenForTest(t)
	persistSession(t, queries, jti, expiresAt)

	router := echo.New()
	router.GET("/", func(c *echo.Context) error {
		service.ClearSession(c)
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderCookie, SessionCookieName+"="+jti)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	cleared := false
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == SessionCookieName && cookie.MaxAge < 0 {
			cleared = true
		}
	}
	if !cleared {
		t.Fatal("expected the session cookie to be expired")
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
		return service.core.Sessions.Start(c.Response(), c.Request(), "user-subject", time.Now().Add(time.Hour))
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

// Rotating cookie flags or the session TTL must not require a restart: the
// cookie helpers read the reloaded config at call time.
func TestSessionCookie_FollowsConfigReload(t *testing.T) {
	loadConfig(t, true)

	if got := (&Auth{core: newCore(nil)}).SessionTTL(); got != 24*time.Hour {
		t.Fatalf("expected the hydrated default TTL, got %v", got)
	}

	expiresAt := time.Now().Add(time.Hour)
	cookie := CookieOptions().Cookie(SessionCookieName, "token", 0, expiresAt)

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

	if cookie := CookieOptions().Cookie(SessionCookieName, "token", 0, expiresAt); !cookie.Secure {
		t.Fatal("expected the session cookie to honour the reloaded cookie_secure")
	}
}

func TestSessionCookie_ExpiresImmediatelyWhenMaxAgeNegative(t *testing.T) {
	if cookie := CookieOptions().Cookie(SessionCookieName, "", -1, time.Time{}); cookie.MaxAge != -1 {
		t.Fatalf("expected MaxAge -1 to delete the cookie, got %d", cookie.MaxAge)
	}
}

func TestRandomToken_ReturnsDistinctHexValues(t *testing.T) {
	first, err := gosso.RandomToken(sessionTokenSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := gosso.RandomToken(sessionTokenSize)
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

	token, err := gosso.RandomToken(sessionTokenSize)
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
		configCore := &Auth{core: &gosso.Auth{Settings: func() gosso.Config { return gosso.Config{} }, CallbackPath: "/api/auth/callback"}}
		var got string
		router.GET("/api/auth/login", func(c *echo.Context) error {
			got = configCore.core.CallbackURL(c.Request())
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, tc.target, nil)
		// gosso.Scheme trusts X-Forwarded-Proto unconditionally: proxies must be in
		// front, or app_url has to be set so the derived scheme is never consulted.
		if tc.proto != "" {
			req.Header.Set("X-Forwarded-Proto", tc.proto)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		if got != tc.want {
			t.Errorf("%s: CallbackURL() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// Behind proxies that do not forward the original scheme (e.g. NetBird), the
// registered redirect URI would be built with http and rejected by the
// provider. app_url overrides the request-derived value in that case.
func TestCallbackURL_PrefersAppUrl(t *testing.T) {
	v := baseConfig()
	v.Set("server.app_url", "https://gocron.example.com")
	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
	t.Cleanup(func() { config.ValidateAndLoadConfig(baseConfig()) })

	router := echo.New()
	configCore := &Auth{core: &gosso.Auth{Settings: settings, CallbackPath: "/api/auth/callback"}}
	var got string
	router.GET("/api/auth/login", func(c *echo.Context) error {
		got = configCore.core.CallbackURL(c.Request())
		return c.NoContent(http.StatusOK)
	})
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "http://gocron:8156/api/auth/login", nil))

	if got != "https://gocron.example.com/api/auth/callback" {
		t.Errorf("CallbackURL() = %q, want %q", got, "https://gocron.example.com/api/auth/callback")
	}
}

// The optional handoff target has to survive the provider round trip inside the
// state cookie and be rejected unless it is a site-relative path.
func TestStartLogin_CarriesValidatedRedirectTarget(t *testing.T) {
	oidcServer := &oidctest.Server{}
	server := httptest.NewServer(oidcServer)
	defer server.Close()
	oidcServer.SetIssuer(server.URL)

	loadConfigWithIssuer(t, server.URL)

	db := newTestDB(t)
	queries := jobs.New(db)

	service := &Auth{sessions: queries, core: newCore(queries)}
	defer service.Shutdown()

	router := echo.New()
	router.GET("/api/auth/login", func(c *echo.Context) error {
		return service.StartLogin(c)
	})

	requests := map[string]string{
		"/api/auth/login?redirect=/jobs/abc":        "/jobs/abc",
		"/api/auth/login?redirect=https://evil.tld": "/",
		"/api/auth/login?redirect=//evil.tld":       "/",
	}

	for target, wantRedirect := range requests {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("%s: expected a redirect to the provider, got %d", target, rec.Code)
		}

		var stateCookie *http.Cookie
		for _, cookie := range rec.Result().Cookies() {
			if strings.HasSuffix(cookie.Name, "_oidc_state") {
				stateCookie = cookie
			}
		}
		if stateCookie == nil {
			t.Fatalf("%s: expected the state cookie", target)
		}

		parts := splitStateValue(stateCookie.Value)
		if len(parts) < 2 {
			t.Fatalf("%s: expected state|verifier in the cookie", target)
		}
		redirect := "/"
		if len(parts) > 2 {
			redirect = parts[2]
		}
		if redirect != wantRedirect {
			t.Errorf("%s: stored redirect %q, want %q", target, redirect, wantRedirect)
		}
	}
}

func splitStateValue(value string) []string {
	return strings.Split(value, "|")
}

// The handoff target stored during StartLogin has to come back out of
// ConsumeLoginState, otherwise the callback cannot honour the deep link and
// every login silently lands on the app root.
func TestConsumeLoginState_ReturnsStoredRedirectTarget(t *testing.T) {
	oidcServer := &oidctest.Server{}
	server := httptest.NewServer(oidcServer)
	defer server.Close()
	oidcServer.SetIssuer(server.URL)

	loadConfigWithIssuer(t, server.URL)

	service := &Auth{core: newCore(nil)}

	router := echo.New()
	router.GET("/api/auth/login", func(c *echo.Context) error {
		return service.StartLogin(c)
	})
	router.GET("/api/auth/callback", func(c *echo.Context) error {
		state, verifier, redirect, err := service.core.ConsumeLoginState(c.Response(), c.Request())
		if err != nil {
			t.Fatalf("consume failed: %v", err)
		}

		if state == "" || verifier == "" {
			t.Fatalf("expected state and verifier, got %q and %q", state, verifier)
		}
		if redirect != "/jobs/abc" {
			t.Fatalf("expected the stored redirect target, got %q", redirect)
		}

		return c.NoContent(http.StatusOK)
	})

	start := httptest.NewRecorder()
	router.ServeHTTP(start, httptest.NewRequest(http.MethodGet, "/api/auth/login?redirect=/jobs/abc", nil))

	if start.Code != http.StatusFound {
		t.Fatalf("expected the login to redirect to the provider, got %d", start.Code)
	}

	// StartLogin clears the stale state cookie before writing the new one, so
	// replay only the live value the browser would end up sending.
	callback := httptest.NewRequest(http.MethodGet, "/api/auth/callback?state=any&code=any", nil)
	for _, cookie := range start.Result().Cookies() {
		if strings.HasSuffix(cookie.Name, "_oidc_state") && cookie.Value != "" {
			callback.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
		}
	}

	router.ServeHTTP(httptest.NewRecorder(), callback)
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
