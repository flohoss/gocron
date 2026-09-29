package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/flohoss/gocron/internal/auth"
)

type mockAuthService struct {
	enabled         bool
	startLoginErr   error
	consumeState    string
	consumeVerifier string
	consumeErr      error
	subject         string
	exchangeErr     error
	startSessionErr error
	authenticateErr error

	startLoginCalled bool
	sessionSubject   string
}

func (m *mockAuthService) Enabled() bool { return m.enabled }

func (m *mockAuthService) StartLogin(c *echo.Context) error {
	m.startLoginCalled = true
	if m.startLoginErr != nil {
		return m.startLoginErr
	}
	return c.Redirect(http.StatusFound, "https://sso.example.com/authorize?state=abc")
}

func (m *mockAuthService) ConsumeLoginState(c *echo.Context) (string, string, error) {
	if m.consumeErr != nil {
		return "", "", m.consumeErr
	}
	return m.consumeState, m.consumeVerifier, nil
}

func (m *mockAuthService) Authenticate(c *echo.Context) (*auth.User, error) {
	if m.authenticateErr != nil {
		return nil, m.authenticateErr
	}
	return &auth.User{Subject: m.subject}, nil
}

func (m *mockAuthService) Exchange(c *echo.Context, code, verifier string) (string, error) {
	if m.exchangeErr != nil {
		return "", m.exchangeErr
	}
	return m.subject, nil
}

func (m *mockAuthService) StartSession(c *echo.Context, subject string, expiresAt time.Time) error {
	if m.startSessionErr != nil {
		return m.startSessionErr
	}
	m.sessionSubject = subject
	return nil
}

func (m *mockAuthService) SessionTTL() time.Duration { return time.Hour }

func (m *mockAuthService) ClearSessionCookie() *http.Cookie {
	return &http.Cookie{Name: auth.SessionCookieName, MaxAge: -1}
}

func (m *mockAuthService) Logout(c *echo.Context) error { return nil }

func (m *mockAuthService) LogoutURL() string { return "" }

func (m *mockAuthService) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}

// runHandler exercises a huma auth handler through a real echo request so that
// cookies and headers land on the recorded response.
func runHandler[O any](t *testing.T, mock *mockAuthService, method string, call func(*AuthHandler, context.Context) (*O, error)) (*httptest.ResponseRecorder, *O) {
	t.Helper()

	var out *O
	var handlerErr error
	router := echo.New()
	router.Add(method, "/probe", func(c *echo.Context) error {
		ah := NewAuthHandler(mock)
		out, handlerErr = call(ah, context.WithValue(c.Request().Context(), echoContextKey{}, c))
		return nil
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(method, "/probe", nil))

	if handlerErr != nil {
		t.Fatalf("handler error: %v", handlerErr)
	}

	return rec, out
}

func cookie(rec *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func serveEchoHandler(t *testing.T, handler echo.HandlerFunc, method, path, target string) *httptest.ResponseRecorder {
	t.Helper()

	router := echo.New()
	router.Add(method, path, handler)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(method, target, nil))

	return rec
}

// The login route must be a top-level redirect: EventSource and the browser
// cannot follow a JSON login handshake, and the state cookie has to survive the
// navigation to the provider.
func TestLoginHandler_RedirectsToProvider(t *testing.T) {
	mock := &mockAuthService{enabled: true}

	rec := serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).loginHandler(c)
	}, http.MethodGet, "/probe", "/probe")

	if !mock.startLoginCalled {
		t.Fatal("expected login flow to be started")
	}
	if rec.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
	if location := rec.Header().Get(echo.HeaderLocation); location == "" {
		t.Fatal("expected a redirect target")
	}
}

func TestLoginHandler_ServerErrorWhenStartLoginFails(t *testing.T) {
	mock := &mockAuthService{enabled: true, startLoginErr: errors.New("no entropy")}

	rec := serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).loginHandler(c)
	}, http.MethodGet, "/probe", "/probe")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

// A provider that cannot be reached is a temporary outage, not a client error:
// the status has to say "try again later" so a proxy or the user does not treat
// it as a bad login.
func TestLoginHandler_ServiceUnavailableWhenProviderUnreachable(t *testing.T) {
	mock := &mockAuthService{enabled: true, startLoginErr: auth.ErrProviderUnavailable}

	rec := serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).loginHandler(c)
	}, http.MethodGet, "/probe", "/probe")

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestLoginHandler_NotFoundWhenAuthDisabled(t *testing.T) {
	mock := &mockAuthService{enabled: false}

	rec := serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).loginHandler(c)
	}, http.MethodGet, "/probe", "/probe")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	if mock.startLoginCalled {
		t.Fatal("expected no login attempt while auth is disabled")
	}
}

func TestCallbackHandler_NotFoundWhenAuthDisabled(t *testing.T) {
	rec := runEchoCallback(t, &mockAuthService{enabled: false}, "state=abc&code=abc")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func runEchoCallback(t *testing.T, mock *mockAuthService, rawQuery string) *httptest.ResponseRecorder {
	t.Helper()

	target := "/probe"
	if rawQuery != "" {
		target += "?" + rawQuery
	}

	return serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).callbackHandler(c)
	}, http.MethodGet, "/probe", target)
}

func TestCallbackHandler_ProviderErrorReturnsJSON(t *testing.T) {
	rec := runEchoCallback(t, &mockAuthService{enabled: true}, "error=access_denied")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"sso_failed", "Authentication required", "access_denied"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %q, got: %s", want, body)
		}
	}
}

func TestCallbackHandler_RejectsStateMismatch(t *testing.T) {
	rec := runEchoCallback(t, &mockAuthService{enabled: true, consumeState: "expected"}, "state=other&code=abc")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"state", "Authentication required"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %q, got: %s", want, body)
		}
	}
}

func TestCallbackHandler_RejectsMissingStateCookie(t *testing.T) {
	rec := runEchoCallback(t, &mockAuthService{enabled: true, consumeErr: auth.ErrUnauthenticated}, "state=abc&code=abc")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

// A successful callback must start the session and land the user on the app root.
func TestCallbackHandler_StartsSessionAndRedirectsHome(t *testing.T) {
	mock := &mockAuthService{
		enabled:      true,
		consumeState: "state-value",
		subject:      "user-subject",
	}

	rec := runEchoCallback(t, mock, "state=state-value&code=abc")

	if location := rec.Header().Get(echo.HeaderLocation); location != "/" {
		t.Fatalf("expected redirect to app root, got %q", location)
	}
	if mock.sessionSubject != "user-subject" {
		t.Fatalf("expected a session for the exchanged subject, got %q", mock.sessionSubject)
	}
}

func TestCallbackHandler_ReturnsJSONWhenExchangeFails(t *testing.T) {
	mock := &mockAuthService{
		enabled:      true,
		consumeState: "state-value",
		exchangeErr:  errors.New("provider down"),
	}

	rec := runEchoCallback(t, mock, "state=state-value&code=abc")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"sso_failed", "Authentication required"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %q, got: %s", want, body)
		}
	}
	if mock.sessionSubject != "" {
		t.Fatal("expected no session after a failed exchange")
	}
}

func TestCallbackHandler_ServiceUnavailableWhenProviderUnreachable(t *testing.T) {
	mock := &mockAuthService{
		enabled:      true,
		consumeState: "state-value",
		exchangeErr:  auth.ErrProviderUnavailable,
	}

	rec := runEchoCallback(t, mock, "state=state-value&code=abc")

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	if mock.sessionSubject != "" {
		t.Fatal("expected no session while the provider is unreachable")
	}
}

func TestCallbackHandler_ServerErrorWhenSessionStartFails(t *testing.T) {
	mock := &mockAuthService{
		enabled:         true,
		consumeState:    "state-value",
		subject:         "user-subject",
		startSessionErr: errors.New("db down"),
	}

	rec := runEchoCallback(t, mock, "state=state-value&code=abc")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestLogoutHandler_ClearsSessionCookie(t *testing.T) {
	rec, _ := runHandler(t, &mockAuthService{enabled: true}, http.MethodPost, func(ah *AuthHandler, ctx context.Context) (*logoutResponse, error) {
		return ah.logoutHandler(ctx, nil)
	})

	cleared := cookie(rec, auth.SessionCookieName)
	if cleared == nil || cleared.MaxAge != -1 {
		t.Fatalf("expected session cookie to be cleared, got %+v", cleared)
	}
}

func TestMeHandler_ReturnsAuthenticatedUser(t *testing.T) {
	_, response := runHandler(t, &mockAuthService{
		enabled: true,
		subject: "user-subject",
	}, http.MethodGet, func(ah *AuthHandler, ctx context.Context) (*currentUserResponse, error) {
		return ah.meHandler(ctx, nil)
	})

	if !response.Body.Authenticated || !response.Body.AuthEnabled {
		t.Fatalf("unexpected response body: %+v", response.Body)
	}
}

// With SSO on but no session, /me reports an anonymous visitor instead of a 401:
// the SPA needs auth_enabled to decide between the login page and the app.
func TestMeHandler_ReportsAnonymousWhenEnabled(t *testing.T) {
	_, response := runHandler(t, &mockAuthService{
		enabled:         true,
		authenticateErr: auth.ErrUnauthenticated,
	}, http.MethodGet, func(ah *AuthHandler, ctx context.Context) (*currentUserResponse, error) {
		return ah.meHandler(ctx, nil)
	})

	if !response.Body.AuthEnabled || response.Body.Authenticated {
		t.Fatalf("unexpected response body: %+v", response.Body)
	}
}

// With SSO off the SPA must be told authentication is not in play, rather than
// being handed a 401 it cannot recover from.
func TestMeHandler_ReportsDisabledAuth(t *testing.T) {
	_, response := runHandler(t, &mockAuthService{enabled: false}, http.MethodGet, func(ah *AuthHandler, ctx context.Context) (*currentUserResponse, error) {
		return ah.meHandler(ctx, nil)
	})

	if response.Body.AuthEnabled || response.Body.Authenticated {
		t.Fatalf("expected disabled auth response, got: %+v", response.Body)
	}
}
