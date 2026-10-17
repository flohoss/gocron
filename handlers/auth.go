package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v5"

	"github.com/flohoss/gocron/internal/auth"
	"github.com/flohoss/gosso/v2"
)

const sessionScheme = "cookieAuth"

type AuthService interface {
	Enabled() bool
	StartLogin(c *echo.Context) error
	HandleCallback(c *echo.Context) (*gosso.LoginResult, *gosso.LoginFailure)
	Authenticate(c *echo.Context) (*auth.User, error)
	ClearSession(c *echo.Context)
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
		slog.Error("Failed to start login", "error", err)
		return redirect(c, LoginErrorPath(string(gosso.ClassifyLoginError(err))))
	}

	return nil
}

func (ah *AuthHandler) callbackHandler(c *echo.Context) error {
	if !ah.Auth.Enabled() {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	result, failure := ah.Auth.HandleCallback(c)
	if failure != nil {
		gosso.LogFailure(failure)

		return redirect(c, LoginErrorPath(string(failure.Code)))
	}

	slog.Info("User logged in", "subject", result.Subject)

	return redirect(c, result.Redirect)
}

type logoutResponse struct {
	Body logoutBody
}

type logoutBody struct {
	Redirect string `json:"redirect" doc:"Where the SPA should navigate after the session was cleared. Always '/'; provider-side logout is out of scope."`
}

func (ah *AuthHandler) logoutHandler(ctx context.Context, input *struct{}) (*logoutResponse, error) {
	c := echoContextFrom(ctx)

	ah.Auth.ClearSession(c)

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
