package handlers

import (
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
	"github.com/flohoss/gocron/services/jobs"
	"github.com/flohoss/gosso/v3"
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

func setupRouterWithAuth(t *testing.T) (*echo.Echo, *gosso.Auth, *sql.DB) {
	t.Helper()

	// The SPA fallback parses web/index.html relative to the working directory,
	// which is the package dir under go test but the repo root at runtime.
	chdirRepoRoot(t)

	loadEnabledAuthConfig(t)

	db := openTestDB(t)

	auth, err := gosso.New(gosso.Options{
		Enabled:      true,
		IssuerURL:    "https://sso.example.com",
		ClientID:     "gocron",
		ClientSecret: "test-secret",
		DB:           db,
	})
	if err != nil {
		t.Fatalf("failed to create auth: %v", err)
	}

	jh := &JobHandler{JobService: &stubJobService{queries: jobs.New(db)}}
	ch := &CommandHandler{CommandsService: &stubCommandsService{}}

	e := InitRouter()
	SetupRouter(e, auth, jh, ch)

	return e, auth, db
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "auth.sqlite"))
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return db
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

	auth, err := gosso.New(gosso.Options{DB: openTestDB(t)})
	if err != nil {
		t.Fatalf("failed to create auth: %v", err)
	}

	jh := &JobHandler{JobService: &stubJobService{}}
	ch := &CommandHandler{CommandsService: &stubCommandsService{}}

	e := InitRouter()
	SetupRouter(e, auth, jh, ch)

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
	e, _, _ := setupRouterWithAuth(t)

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
		{method: http.MethodGet, path: "/api/auth/callback?code=x&state=y", want: http.StatusFound},
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
	e, _, _ := setupRouterWithAuth(t)

	for _, path := range []string{"/assets/app.js", "/static/site.webmanifest"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

		if rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound {
			t.Errorf("%s = %d, expected the asset route to bypass the session check", path, rec.Code)
		}
	}
}

func TestSetupRouter_WithAuth_RedirectsShellToProvider(t *testing.T) {
	e, _, _ := setupRouterWithAuth(t)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if location := rec.Header().Get(echo.HeaderLocation); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

// A deep link must land the user on the app once the session is established.
func TestSetupRouter_WithAuth_ServesShellWithValidSession(t *testing.T) {
	e, auth, db := setupRouterWithAuth(t)

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		t.Fatalf("failed to generate session token: %v", err)
	}
	jti := hex.EncodeToString(token)
	if _, err := db.ExecContext(t.Context(),
		`INSERT INTO sso_sessions (jti, subject, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		jti, "user-subject", time.Now().UnixMilli(), time.Now().Add(time.Hour).UnixMilli()); err != nil {
		t.Fatalf("failed to persist session: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/jobs/some-job", nil)
	req.Header.Set(echo.HeaderCookie, auth.SessionCookieName()+"="+jti)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected the shell to be served with a valid session, got %d", rec.Code)
	}
}

// The generated frontend client and the docs page need the session cookie scheme
// to be part of the published spec.
func TestSetupRouter_WithAuth_DocumentsSessionScheme(t *testing.T) {
	e, auth, _ := setupRouterWithAuth(t)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))

	spec := rec.Body.String()
	for _, expected := range []string{`"cookieAuth"`, `"in":"cookie"`, auth.SessionCookieName()} {
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
