package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
)

const (
	SessionCookieName = "gocron_session"
	stateCookieName   = "gocron_oidc_state"

	callbackPath = "/api/auth/callback"

	stateCookieTTL = 10 * time.Minute
	stateSeparator = "|"

	sessionTokenSize = 32
)

var (
	ErrUnauthenticated     = errors.New("not authenticated")
	ErrProviderUnavailable = errors.New("identity provider is not reachable")

	scopes        = []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail}
	usernameOrder = []string{"preferred_username", "name", "email"}
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type providerState struct {
	key           string
	verifier      *oidc.IDTokenVerifier
	oauthConfig   oauth2.Config
	endSessionURL string
}

type Auth struct {
	providerMu sync.Mutex
	provider   *providerState

	sessions SessionStore
	stop     context.CancelFunc
	wg       sync.WaitGroup
}

type SessionStore interface {
	CreateSession(ctx context.Context, arg jobs.CreateSessionParams) error
	GetActiveSession(ctx context.Context, arg jobs.GetActiveSessionParams) (jobs.Session, error)
	RevokeSession(ctx context.Context, jti string) error
	DeleteExpiredSessions(ctx context.Context, expiresAt int64) error
}

func New(sessions SessionStore) *Auth {
	auth := &Auth{sessions: sessions}
	auth.startCleanup()
	auth.warmProvider()

	return auth
}

func (a *Auth) warmProvider() {
	if !a.Enabled() {
		return
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		_, _ = a.resolve()
	}()
}

func (a *Auth) Enabled() bool {
	return a != nil && a.settings().Enabled
}

func (a *Auth) settings() config.OIDCSettings {
	return config.GetAuth().OIDC
}

func (a *Auth) resolve() (*providerState, error) {
	settings := a.settings()
	key := providerKey(settings)

	a.providerMu.Lock()
	defer a.providerMu.Unlock()

	if a.provider != nil && a.provider.key == key {
		return a.provider, nil
	}

	provider, err := discoverProvider(settings)
	if err != nil {
		return nil, err
	}

	a.provider = provider
	return provider, nil
}

func providerKey(settings config.OIDCSettings) string {
	return settings.IssuerURL + "|" + settings.ClientID
}

func discoverProvider(settings config.OIDCSettings) (*providerState, error) {
	provider, err := oidc.NewProvider(context.Background(), settings.IssuerURL)
	if err != nil {
		slog.Warn("OIDC discovery failed, will retry on the next request", "issuer", settings.IssuerURL, "error", err)
		return nil, fmt.Errorf("%w: discovery for %q failed: %w", ErrProviderUnavailable, settings.IssuerURL, err)
	}

	var claims struct {
		EndSessionURL string `json:"end_session_endpoint"`
	}
	if err := provider.Claims(&claims); err != nil {
		slog.Warn("OIDC provider metadata could not be read, will retry on the next request", "issuer", settings.IssuerURL, "error", err)
		return nil, fmt.Errorf("%w: reading %q metadata failed: %w", ErrProviderUnavailable, settings.IssuerURL, err)
	}

	slog.Info("OIDC provider discovered", "issuer", settings.IssuerURL)

	return &providerState{
		key:           providerKey(settings),
		verifier:      provider.Verifier(&oidc.Config{ClientID: settings.ClientID}),
		oauthConfig:   oauth2.Config{ClientID: settings.ClientID, Endpoint: provider.Endpoint(), Scopes: scopes},
		endSessionURL: claims.EndSessionURL,
	}, nil
}

func (a *Auth) Shutdown() {
	if a.stop != nil {
		a.stop()
	}
	a.wg.Wait()
}

func (a *Auth) startCleanup() {
	ctx, cancel := context.WithCancel(context.Background())
	a.stop = cancel

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		a.deleteExpiredSessions()

		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				a.deleteExpiredSessions()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a *Auth) deleteExpiredSessions() {
	if a.sessions == nil {
		return
	}

	if err := a.sessions.DeleteExpiredSessions(context.Background(), time.Now().UnixMilli()); err != nil {
		slog.Warn("Failed to delete expired sessions", "error", err)
	}
}

func (a *Auth) StartLogin(c *echo.Context) error {
	provider, err := a.resolve()
	if err != nil {
		return err
	}

	state, err := randomToken(24)
	if err != nil {
		return err
	}

	verifier := oauth2.GenerateVerifier()

	a.clearLoginState(c)
	c.SetCookie(a.cookie(stateCookieName, state+stateSeparator+verifier, int(stateCookieTTL.Seconds()), time.Time{}))

	login := a.loginConfig(provider, c)
	c.Redirect(http.StatusFound, login.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)))
	return nil
}

func (a *Auth) loginConfig(provider *providerState, c *echo.Context) oauth2.Config {
	login := provider.oauthConfig
	login.ClientSecret = a.settings().ClientSecret
	login.RedirectURL = callbackURL(c)
	return login
}

func (a *Auth) ConsumeLoginState(c *echo.Context) (string, string, error) {
	a.clearLoginState(c)

	cookie, err := c.Cookie(stateCookieName)
	if err != nil {
		return "", "", ErrUnauthenticated
	}

	state, verifier, found := strings.Cut(cookie.Value, stateSeparator)
	if !found || state == "" || verifier == "" {
		return "", "", ErrUnauthenticated
	}

	return state, verifier, nil
}

func (a *Auth) clearLoginState(c *echo.Context) {
	c.SetCookie(a.cookie(stateCookieName, "", -1, time.Time{}))
}

func callbackURL(c *echo.Context) string {
	return c.Scheme() + "://" + c.Request().Host + callbackPath
}

func (a *Auth) Exchange(c *echo.Context, code, verifier string) (*User, error) {
	provider, err := a.resolve()
	if err != nil {
		return nil, err
	}

	login := a.loginConfig(provider, c)
	token, err := login.Exchange(c.Request().Context(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, errors.New("OIDC provider did not return an id_token")
	}

	idToken, err := provider.verifier.Verify(c.Request().Context(), rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id_token: %w", err)
	}

	var claims map[string]any
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to decode id_token claims: %w", err)
	}

	if !isEmailVerified(claims["email_verified"]) {
		return nil, fmt.Errorf("OIDC provider returned an unverified email claim for %q", idToken.Subject)
	}

	email, _ := claims["email"].(string)
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("OIDC provider returned no usable email claim for %q", idToken.Subject)
	}

	return &User{Email: email, Username: a.username(claims)}, nil
}

func isEmailVerified(claim any) bool {
	switch value := claim.(type) {
	case bool:
		return value
	case string:
		parsed, err := strconv.ParseBool(value)
		return err == nil && parsed
	default:
		return false
	}
}

func (a *Auth) username(claims map[string]any) string {
	for _, claim := range usernameOrder {
		if value, _ := claims[claim].(string); value != "" {
			return value
		}
	}

	return ""
}

func (a *Auth) StartSession(c *echo.Context, user *User, expiresAt time.Time) error {
	jti, err := randomToken(sessionTokenSize)
	if err != nil {
		return err
	}

	if a.sessions != nil {
		if err := a.sessions.CreateSession(c.Request().Context(), jobs.CreateSessionParams{
			Jti:       jti,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: time.Now().UnixMilli(),
			ExpiresAt: expiresAt.UnixMilli(),
		}); err != nil {
			return fmt.Errorf("failed to persist session: %w", err)
		}
	}

	c.SetCookie(a.cookie(SessionCookieName, jti, 0, expiresAt))
	return nil
}

func (a *Auth) cookie(name, value string, maxAge int, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expires,
		HttpOnly: true,
		Secure:   a.settings().CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (a *Auth) SessionCookie(value string, expiresAt time.Time) *http.Cookie {
	return a.cookie(SessionCookieName, value, 0, expiresAt)
}

func (a *Auth) ClearSessionCookie() *http.Cookie {
	return a.cookie(SessionCookieName, "", -1, time.Time{})
}

func (a *Auth) SessionTTL() time.Duration {
	ttl := a.settings().SessionTTL
	if ttl == 0 {
		ttl = 24 * time.Hour
	}
	return ttl
}

func (a *Auth) sessionToken(c *echo.Context) (string, error) {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", ErrUnauthenticated
	}

	return cookie.Value, nil
}

func (a *Auth) Authenticate(c *echo.Context) (*User, error) {
	token, err := a.sessionToken(c)
	if err != nil {
		return nil, err
	}

	if a.sessions == nil {
		return nil, ErrUnauthenticated
	}

	session, err := a.sessions.GetActiveSession(c.Request().Context(), jobs.GetActiveSessionParams{
		Jti:       token,
		ExpiresAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, ErrUnauthenticated
	}

	return &User{Email: session.Email, Username: session.Username}, nil
}

func (a *Auth) Logout(c *echo.Context) error {
	token, err := a.sessionToken(c)
	if err != nil || a.sessions == nil {
		return nil
	}

	return a.sessions.RevokeSession(c.Request().Context(), token)
}

func (a *Auth) EndSessionURL(redirectTo string) string {
	provider, err := a.resolve()
	if err != nil {
		return ""
	}

	endpoint, err := url.Parse(provider.endSessionURL)
	if err != nil {
		return ""
	}

	query := endpoint.Query()
	query.Set("post_logout_redirect_uri", redirectTo)
	endpoint.RawQuery = query.Encode()

	return endpoint.String()
}

func (a *Auth) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if !a.Enabled() {
				return next(c)
			}

			if _, err := a.Authenticate(c); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			return next(c)
		}
	}
}

func randomToken(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), nil
}
