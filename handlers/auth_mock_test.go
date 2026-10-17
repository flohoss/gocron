package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/flohoss/gocron/internal/auth"
	"github.com/flohoss/gosso/v2"
)

// mockAuthService is a scriptable AuthService. Every failure knob returns its
// error so the handler branches under test stay reachable.
type mockAuthService struct {
	enabled         bool
	startLoginErr   error
	callbackResult  *gosso.LoginResult
	callbackFailure *gosso.LoginFailure
	authenticateErr error

	startLoginCalled bool
	sessionSubject   string
	sessionCleared   bool
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

func (m *mockAuthService) HandleCallback(c *echo.Context) (*gosso.LoginResult, *gosso.LoginFailure) {
	if m.callbackFailure != nil {
		return nil, m.callbackFailure
	}

	m.sessionSubject = m.callbackResult.Subject

	return m.callbackResult, nil
}

func (m *mockAuthService) Authenticate(c *echo.Context) (*auth.User, error) {
	if m.authenticateErr != nil {
		return nil, m.authenticateErr
	}

	return &auth.User{Subject: "user-subject"}, nil
}

func (m *mockAuthService) SessionTTL() time.Duration {
	return time.Hour
}

func (m *mockAuthService) ClearSession(c *echo.Context) {
	m.sessionCleared = true
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
