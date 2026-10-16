package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v5"

	"github.com/flohoss/gocron/internal/auth"
	"github.com/flohoss/gosso"
)

const sessionScheme = "cookieAuth"

type AuthService interface {
	Enabled() bool
	StartLogin(c *echo.Context) error
	ConsumeLoginState(c *echo.Context) (string, string, string, error)
	Authenticate(c *echo.Context) (*auth.User, error)
	Exchange(c *echo.Context, code, verifier string) (string, error)
	StartSession(c *echo.Context, subject string, expiresAt time.Time) error
	SessionTTL() time.Duration
	ClearSessionCookie() *http.Cookie
	Logout(c *echo.Context) error
	Middleware() echo.MiddlewareFunc
}

type AuthHandler struct {
	Auth AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{Auth: service}
}

func (ah *AuthHandler) Register(api huma.API) {
	huma.Register(api, ah.logoutOperation(), ah.logoutHandler)
	huma.Register(api, ah.meOperation(), ah.meHandler)
}

func (ah *AuthHandler) logoutOperation() huma.Operation {
	return huma.Operation{
		OperationID: "post-auth-logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Summary:     "Logout",
		Description: "Revokes the current session and clears the session cookie.",
		Tags:        []string{"Auth"},
		Security:    []map[string][]string{{sessionScheme: {}}},
	}
}

func (ah *AuthHandler) meOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-auth-me",
		Method:      http.MethodGet,
		Path:        "/api/auth/me",
		Summary:     "Get current user",
		Description: "Returns whether single sign-on is enabled and the identity of the current session.",
		Tags:        []string{"Auth"},
		Security:    []map[string][]string{{sessionScheme: {}}},
	}
}

type currentUserBody struct {
	Authenticated bool `json:"authenticated"`
	AuthEnabled   bool `json:"auth_enabled" doc:"Whether single sign-on is enabled on this instance."`
}

type currentUserResponse struct {
	Body currentUserBody
}

func (ah *AuthHandler) loginHandler(c *echo.Context) error {
	if !ah.Auth.Enabled() {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	if err := ah.Auth.StartLogin(c); err != nil {
		if errors.Is(err, gosso.ErrProviderUnavailable) {
			slog.Error("OIDC provider unavailable", "error", err)
			return redirect(c, LoginErrorPath(string(gosso.LoginErrorUnavailable)))
		}

		slog.Error("Failed to start login", "error", err)
		return redirect(c, LoginErrorPath(string(gosso.LoginErrorFailed)))
	}

	return nil
}

func (ah *AuthHandler) callbackHandler(c *echo.Context) error {
	if !ah.Auth.Enabled() {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	if providerError := c.QueryParam("error"); providerError != "" {
		slog.Warn("OIDC provider rejected login", "error", providerError)
		return redirect(c, LoginErrorPath(string(gosso.LoginErrorRejected)))
	}

	state, verifier, target, err := ah.Auth.ConsumeLoginState(c)
	if err != nil {
		return redirect(c, LoginErrorPath(string(gosso.LoginErrorExpired)))
	}

	if c.QueryParam("state") != state {
		slog.Warn("OIDC state mismatch, discarding callback")
		return redirect(c, LoginErrorPath(string(gosso.LoginErrorExpired)))
	}

	subject, err := ah.Auth.Exchange(c, c.QueryParam("code"), verifier)
	if err != nil {
		if errors.Is(err, gosso.ErrProviderUnavailable) {
			slog.Error("OIDC provider unavailable", "error", err)
			return redirect(c, LoginErrorPath(string(gosso.LoginErrorUnavailable)))
		}

		slog.Error("OIDC login failed", "error", err)
		return redirect(c, LoginErrorPath(string(gosso.ClassifyLoginError(err))))
	}

	if err := ah.Auth.StartSession(c, subject, time.Now().Add(ah.Auth.SessionTTL())); err != nil {
		slog.Error("Failed to create session", "error", err)
		return redirect(c, LoginErrorPath(string(gosso.LoginErrorFailed)))
	}

	slog.Info("User logged in", "subject", subject)

	if target == "" {
		target = "/"
	}

	return redirect(c, target)
}

type logoutResponse struct {
	Body logoutBody
}

type logoutBody struct {
	Redirect string `json:"redirect" doc:"Where the SPA should navigate after the session was cleared. Always '/'; provider-side logout is out of scope."`
}

func (ah *AuthHandler) logoutHandler(ctx context.Context, input *struct{}) (*logoutResponse, error) {
	c := echoContextFrom(ctx)

	if err := ah.Auth.Logout(c); err != nil {
		slog.Warn("Failed to revoke session", "error", err)
	}

	c.SetCookie(ah.Auth.ClearSessionCookie())

	return &logoutResponse{Body: logoutBody{Redirect: "/"}}, nil
}

func (ah *AuthHandler) meHandler(ctx context.Context, input *struct{}) (*currentUserResponse, error) {
	c := echoContextFrom(ctx)

	if !ah.Auth.Enabled() {
		return &currentUserResponse{Body: currentUserBody{AuthEnabled: false}}, nil
	}

	if _, err := ah.Auth.Authenticate(c); err != nil {
		return &currentUserResponse{Body: currentUserBody{AuthEnabled: true}}, nil
	}

	return &currentUserResponse{Body: currentUserBody{
		Authenticated: true,
		AuthEnabled:   true,
	}}, nil
}
