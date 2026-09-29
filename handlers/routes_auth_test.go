package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/internal/auth"
	"github.com/flohoss/gocron/services/jobs"
)

func loadEnabledAuthConfig(t *testing.T) {
	t.Helper()

	v := viper.New()
	v.Set("time_zone", "UTC")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("auth.oidc.enabled", true)
	v.Set("auth.oidc.issuer_url", "https://sso.example.com")
	v.Set("auth.oidc.client_id", "gocron")
	v.Set("auth.oidc.client_secret", "test-signing-key")
	v.Set("jobs", []map[string]any{{
		"name":     "Router Auth Test Job",
		"commands": []string{"echo test"},
	}})

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
}

// stubAuthService embeds the real Auth so the router test exercises the actual
// middleware and cookie handling without a live OIDC provider.
type stubAuthService struct {
	*auth.Auth
}

func setupRouterWithAuth(t *testing.T) (*echo.Echo, *jobs.Queries) {
	t.Helper()

	// The SPA fallback parses web/index.html relative to the working directory,
	// which is the package dir under go test but the repo root at runtime.
	chdirRepoRoot(t)

	loadEnabledAuthConfig(t)

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "router.sqlite"))
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		jti TEXT PRIMARY KEY,
		subject TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		expires_at INTEGER NOT NULL,
		revoked INTEGER NOT NULL DEFAULT 0
	)`); err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	queries := jobs.New(db)
	service := auth.New(queries)
	t.Cleanup(service.Shutdown)

	ah := NewAuthHandler(&stubAuthService{service})
	jh := &JobHandler{JobService: &stubJobService{queries: queries}}
	ch := &CommandHandler{CommandsService: &stubCommandsService{}}

	e := InitRouter()
	SetupRouter(e, jh, ch, ah)

	return e, queries
}

func statusOf(e *echo.Echo, method, target string) int {
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec.Code
}

// The OIDC routes always exist so enabling SSO at runtime works without a
// restart, but with auth disabled they must behave like any unknown /api route
// rather than leaking a handshake.
func TestSetupRouter_WithoutAuth_DisablesTheOIDCHandshake(t *testing.T) {
	chdirRepoRoot(t)

	v := viper.New()
	v.Set("time_zone", "UTC")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("jobs", []map[string]any{{
		"name":     "Router Auth Test Job",
		"commands": []string{"echo test"},
	}})
	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	ah := NewAuthHandler(&mockAuthService{enabled: false})
	jh := &JobHandler{JobService: &stubJobService{}}
	ch := &CommandHandler{CommandsService: &stubCommandsService{}}

	e := InitRouter()
	SetupRouter(e, jh, ch, ah)

	for _, path := range []string{"/api/auth/login", "/api/auth/callback?code=x&state=y"} {
		if got := statusOf(e, http.MethodGet, path); got != http.StatusNotFound {
			t.Errorf("%s = %d, want %d", path, got, http.StatusNotFound)
		}
	}
}

// With auth enabled only the data and execution API is protected. The SPA shell
// redirects to /login instead of being served, while assets, the OIDC handshake
// and the API docs stay public.
func TestSetupRouter_WithAuth_ProtectsOnlyTheDataAPI(t *testing.T) {
	e, _ := setupRouterWithAuth(t)

	cases := []struct {
		method string
		path   string
		want   int
	}{
		{method: http.MethodGet, path: "/api/jobs", want: http.StatusUnauthorized},
		{method: http.MethodPost, path: "/api/command", want: http.StatusUnauthorized},
		{method: http.MethodGet, path: "/api/events?stream=status", want: http.StatusUnauthorized},
		{method: http.MethodGet, path: "/api/openapi.json", want: http.StatusOK},
		{method: http.MethodGet, path: "/api/docs", want: http.StatusOK},
		{method: http.MethodGet, path: "/api/auth/callback?code=x&state=y", want: http.StatusUnauthorized},
		{method: http.MethodGet, path: "/health", want: http.StatusOK},
		{method: http.MethodGet, path: "/robots.txt", want: http.StatusOK},
		{method: http.MethodGet, path: "/", want: http.StatusFound},
		{method: http.MethodGet, path: "/jobs/some-job", want: http.StatusFound},
	}

	for _, tc := range cases {
		if got := statusOf(e, tc.method, tc.path); got != tc.want {
			t.Errorf("%s %s = %d, want %d", tc.method, tc.path, got, tc.want)
		}
	}
}

// Asset routes sit outside the session middleware: they must neither be
// redirected to /login nor answer 401. The exact status depends on whether a
// build is present, so only the bypass is asserted here.
func TestSetupRouter_WithAuth_ServesAssetsWithoutSession(t *testing.T) {
	e, _ := setupRouterWithAuth(t)

	for _, path := range []string{"/assets/app.js", "/static/site.webmanifest"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

		if rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound {
			t.Errorf("%s = %d, expected the asset route to bypass the session check", path, rec.Code)
		}
	}
}

func TestSetupRouter_WithAuth_RedirectsShellToProvider(t *testing.T) {
	e, _ := setupRouterWithAuth(t)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if location := rec.Header().Get(echo.HeaderLocation); location != "/api/auth/login" {
		t.Fatalf("expected redirect to /api/auth/login, got %q", location)
	}
}

// A deep link must land the user on the app once the session is established.
func TestSetupRouter_WithAuth_ServesShellWithValidSession(t *testing.T) {
	e, queries := setupRouterWithAuth(t)

	expiresAt := time.Now().Add(time.Hour)
	jti := make([]byte, 32)
	if _, err := rand.Read(jti); err != nil {
		t.Fatalf("failed to generate session token: %v", err)
	}
	token := hex.EncodeToString(jti)
	if err := queries.CreateSession(context.Background(), jobs.CreateSessionParams{
		Jti:       token,
		Subject:   "user-subject",
		CreatedAt: time.Now().UnixMilli(),
		ExpiresAt: expiresAt.UnixMilli(),
	}); err != nil {
		t.Fatalf("failed to persist session: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/jobs/some-job", nil)
	req.Header.Set(echo.HeaderCookie, auth.SessionCookieName+"="+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected the shell to be served with a valid session, got %d", rec.Code)
	}
}

// The generated frontend client and the docs page need the session cookie scheme
// to be part of the published spec.
func TestSetupRouter_WithAuth_DocumentsSessionScheme(t *testing.T) {
	e, _ := setupRouterWithAuth(t)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))

	spec := rec.Body.String()
	for _, expected := range []string{`"cookieAuth"`, `"in":"cookie"`, auth.SessionCookieName} {
		if !strings.Contains(spec, expected) {
			t.Errorf("expected %q in the OpenAPI spec", expected)
		}
	}
}

func chdirRepoRoot(t *testing.T) {
	t.Helper()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for !fileExists(filepath.Join(original, "go.mod")) {
		parent := filepath.Dir(original)
		if parent == original {
			t.Fatal("failed to locate the repository root")
		}
		original = parent
	}

	if err := os.Chdir(original); err != nil {
		t.Fatalf("failed to chdir to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
