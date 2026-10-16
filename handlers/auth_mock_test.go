package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/flohoss/gocron/internal/auth"
)

// mockAuthService is a scriptable AuthService. Every failure knob returns its
// error so the handler branches under test stay reachable.
type mockAuthService struct {
	enabled         bool
	startLoginErr   error
	consumeState    string
	consumeVerifier string
	consumeRedirect string
	consumeErr      error
	subject         string
	exchangeErr     error
	startSessionErr error
	authenticateErr error

	startLoginCalled bool
	sessionSubject   string
}

func (m *mockAuthService) Enabled() bool {
	return m.enabled
}

func (m *mockAuthService) StartLogin(c *echo.Context) error {
	m.startLoginCalled = true

	if m.startLoginErr != nil {
		return m.startLoginErr
	}

	return c.Redirect(http.StatusFound, "https://sso.example.com/authorize?state=abc")
}

func (m *mockAuthService) ConsumeLoginState(c *echo.Context) (string, string, string, error) {
	if m.consumeErr != nil {
		return "", "", "", m.consumeErr
	}

	return m.consumeState, m.consumeVerifier, m.consumeRedirect, nil
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

func (m *mockAuthService) SessionTTL() time.Duration {
	return time.Hour
}

func (m *mockAuthService) ClearSessionCookie() *http.Cookie {
	return auth.ClearSessionCookie()
}

func (m *mockAuthService) Logout(c *echo.Context) error {
	return nil
}

func (m *mockAuthService) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if !m.Enabled() {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
		}
	}
}

var _ AuthService = (*mockAuthService)(nil)
