package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	"github.com/flohoss/gosso/v2"
)

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

func serveEchoHandler(t *testing.T, handler echo.HandlerFunc, method, path, target string) *httptest.ResponseRecorder {
	t.Helper()

	router := echo.New()
	router.Add(method, path, handler)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(method, target, nil))

	return rec
}

func runLoginHandler(t *testing.T, mock *mockAuthService) *httptest.ResponseRecorder {
	t.Helper()

	return serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).loginHandler(c)
	}, http.MethodGet, "/probe", "/probe")
}

func runCallbackHandler(t *testing.T, mock *mockAuthService, rawQuery string) *httptest.ResponseRecorder {
	t.Helper()

	target := "/probe"
	if rawQuery != "" {
		target += "?" + rawQuery
	}

	return serveEchoHandler(t, func(c *echo.Context) error {
		return (&AuthHandler{Auth: mock}).callbackHandler(c)
	}, http.MethodGet, "/probe", target)
}

func redirectTarget(rec *httptest.ResponseRecorder) string {
	return rec.Header().Get(echo.HeaderLocation)
}

// The login route has to hand the browser off to the provider, not answer JSON:
// EventSource and the top-level navigation both need a redirect.
func TestLoginHandler_RedirectsToProvider(t *testing.T) {
	mock := &mockAuthService{enabled: true}

	rec := runLoginHandler(t, mock)

	if !mock.startLoginCalled {
		t.Fatal("expected login flow to be started")
	}
	if rec.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}

// A provider outage is a retryable condition, not a bad login, so the user is
// sent back to the login page with a code the SPA can explain.
func TestLoginHandler_RedirectsWithUnavailableWhenProviderUnreachable(t *testing.T) {
	mock := &mockAuthService{enabled: true, startLoginErr: gosso.ErrProviderUnavailable}

	rec := runLoginHandler(t, mock)

	if got := redirectTarget(rec); got != "/login?error=unavailable" {
		t.Fatalf("expected the unavailable login error, got %q", got)
	}
}

func TestLoginHandler_RedirectsWithFailedWhenStartLoginFails(t *testing.T) {
	mock := &mockAuthService{enabled: true, startLoginErr: errors.New("no entropy")}

	rec := runLoginHandler(t, mock)

	if got := redirectTarget(rec); got != "/login?error=failed" {
		t.Fatalf("expected the failed login error, got %q", got)
	}
}

func TestLoginHandler_NotFoundWhenAuthDisabled(t *testing.T) {
	mock := &mockAuthService{enabled: false}

	rec := runLoginHandler(t, mock)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	if mock.startLoginCalled {
		t.Fatal("expected no login attempt while auth is disabled")
	}
}

func TestCallbackHandler_NotFoundWhenAuthDisabled(t *testing.T) {
	rec := runCallbackHandler(t, &mockAuthService{enabled: false}, "state=abc&code=abc")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// An error forwarded by the provider means the user declined or the request was
// rejected upstream; surface it on the login page rather than as a dead end.
func TestCallbackHandler_RedirectsWithRejectedWhenProviderErrors(t *testing.T) {
	mock := &mockAuthService{
		enabled:         true,
		callbackFailure: &gosso.LoginFailure{Code: gosso.LoginErrorRejected, Cause: errors.New("provider reported access_denied")},
	}

	rec := runCallbackHandler(t, mock, "error=access_denied")

	if got := redirectTarget(rec); got != "/login?error=rejected" {
		t.Fatalf("expected the rejected login error, got %q", got)
	}
}

func TestCallbackHandler_RedirectsWithExpiredWhenStateCookieMissing(t *testing.T) {
	mock := &mockAuthService{enabled: true, callbackFailure: &gosso.LoginFailure{Code: gosso.LoginErrorExpired, Cause: gosso.ErrUnauthenticated}}

	rec := runCallbackHandler(t, mock, "state=abc&code=abc")

	if got := redirectTarget(rec); got != "/login?error=expired" {
		t.Fatalf("expected the expired login error, got %q", got)
	}
}

func TestCallbackHandler_RedirectsWithExpiredOnStateMismatch(t *testing.T) {
	mock := &mockAuthService{enabled: true, callbackFailure: &gosso.LoginFailure{Code: gosso.LoginErrorExpired, Cause: gosso.ErrStateMismatch}}

	rec := runCallbackHandler(t, mock, "state=other&code=abc")

	if got := redirectTarget(rec); got != "/login?error=expired" {
		t.Fatalf("expected the expired login error, got %q", got)
	}
}

// The page the user originally asked for is carried through the provider round
// trip, so the callback has to land there instead of the app root.
func TestCallbackHandler_StartsSessionAndRedirectsToRequestedPage(t *testing.T) {
	mock := &mockAuthService{
		enabled:        true,
		callbackResult: &gosso.LoginResult{Subject: "user-subject", Redirect: "/jobs/abc"},
	}

	rec := runCallbackHandler(t, mock, "state=state-value&code=abc")

	if got := redirectTarget(rec); got != "/jobs/abc" {
		t.Fatalf("expected redirect to the requested page, got %q", got)
	}
	if mock.sessionSubject != "user-subject" {
		t.Fatalf("expected a session for the exchanged subject, got %q", mock.sessionSubject)
	}
}

func TestCallbackHandler_RedirectsHomeWhenNoRedirectStored(t *testing.T) {
	mock := &mockAuthService{
		enabled:        true,
		callbackResult: &gosso.LoginResult{Subject: "user-subject", Redirect: "/"},
	}

	rec := runCallbackHandler(t, mock, "state=state-value&code=abc")

	if got := redirectTarget(rec); got != "/" {
		t.Fatalf("expected redirect to the app root, got %q", got)
	}
}

func TestCallbackHandler_RedirectsWithUnavailableWhenProviderUnreachable(t *testing.T) {
	mock := &mockAuthService{
		enabled:         true,
		callbackFailure: &gosso.LoginFailure{Code: gosso.LoginErrorUnavailable, Cause: gosso.ErrProviderUnavailable},
	}

	rec := runCallbackHandler(t, mock, "state=state-value&code=abc")

	if got := redirectTarget(rec); got != "/login?error=unavailable" {
		t.Fatalf("expected the unavailable login error, got %q", got)
	}
	if mock.sessionSubject != "" {
		t.Fatal("expected no session while the provider is unreachable")
	}
}

func TestCallbackHandler_RedirectsWithFailedWhenExchangeFails(t *testing.T) {
	mock := &mockAuthService{
		enabled:         true,
		callbackFailure: &gosso.LoginFailure{Code: gosso.LoginErrorFailed, Cause: errors.New("provider down")},
	}

	rec := runCallbackHandler(t, mock, "state=state-value&code=abc")

	if got := redirectTarget(rec); got != "/login?error=failed" {
		t.Fatalf("expected the failed login error, got %q", got)
	}
	if mock.sessionSubject != "" {
		t.Fatal("expected no session after a failed exchange")
	}
}

func TestCallbackHandler_RedirectsWithFailedWhenSessionStartFails(t *testing.T) {
	mock := &mockAuthService{
		enabled:         true,
		callbackFailure: &gosso.LoginFailure{Code: gosso.LoginErrorFailed, Cause: errors.New("db down")},
	}

	rec := runCallbackHandler(t, mock, "state=state-value&code=abc")

	if got := redirectTarget(rec); got != "/login?error=failed" {
		t.Fatalf("expected the failed login error, got %q", got)
	}
}

func TestLogoutHandler_ClearsSession(t *testing.T) {
	mock := &mockAuthService{enabled: true}

	_, response := runHandler(t, mock, http.MethodPost, func(ah *AuthHandler, ctx context.Context) (*logoutResponse, error) {
		return ah.logoutHandler(ctx, nil)
	})

	if !mock.sessionCleared {
		t.Fatal("expected the session to be cleared")
	}
	if response.Body.Redirect != "/" {
		t.Fatalf("expected the SPA to be sent to the app root, got %q", response.Body.Redirect)
	}
}

func TestMeHandler_ReturnsAuthenticatedUser(t *testing.T) {
	_, response := runHandler(t, &mockAuthService{
		enabled:         true,
		authenticateErr: nil,
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
		authenticateErr: gosso.ErrUnauthenticated,
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
