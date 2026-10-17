package auth

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
	"github.com/flohoss/gosso/v2"
)

const (
	SessionCookieName = "gocron_session"

	callbackPath = "/api/auth/callback"

	sessionTokenSize = 32
)

type User = gosso.User

type SessionStore interface {
	CreateSession(ctx context.Context, arg jobs.CreateSessionParams) error
	GetActiveSession(ctx context.Context, arg jobs.GetActiveSessionParams) (jobs.Session, error)
	DeleteSession(ctx context.Context, jti string) error
	DeleteExpiredSessions(ctx context.Context, expiresAt int64) error
}

type DBSessions struct {
	sessions SessionStore
}

func NewDBSessions(sessions SessionStore) *DBSessions {
	return &DBSessions{sessions: sessions}
}

func (d *DBSessions) Start(w http.ResponseWriter, r *http.Request, subject string, expiresAt time.Time) error {
	jti, err := gosso.RandomToken(sessionTokenSize)
	if err != nil {
		return err
	}

	if d.sessions != nil {
		if err := d.sessions.CreateSession(r.Context(), jobs.CreateSessionParams{
			Jti:       jti,
			Subject:   subject,
			CreatedAt: time.Now().UnixMilli(),
			ExpiresAt: expiresAt.UnixMilli(),
		}); err != nil {
			return err
		}
	}

	http.SetCookie(w, CookieOptions().Cookie(SessionCookieName, jti, 0, expiresAt))

	return nil
}

func (d *DBSessions) Authenticate(w http.ResponseWriter, r *http.Request) (*gosso.User, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, gosso.ErrUnauthenticated
	}

	if d.sessions == nil {
		return nil, gosso.ErrUnauthenticated
	}

	session, err := d.sessions.GetActiveSession(r.Context(), jobs.GetActiveSessionParams{
		Jti:       cookie.Value,
		ExpiresAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, gosso.ErrUnauthenticated
	}

	return &gosso.User{Subject: session.Subject}, nil
}

func (d *DBSessions) Clear(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(SessionCookieName)
	if err == nil && cookie.Value != "" && d.sessions != nil {
		if err := d.sessions.DeleteSession(r.Context(), cookie.Value); err != nil {
			slog.Warn("Failed to revoke session", "error", err)
		}
	}

	http.SetCookie(w, CookieOptions().Cookie(SessionCookieName, "", -1, time.Time{}))

	return nil
}

func (d *DBSessions) CookieName() string {
	return SessionCookieName
}

func CookieOptions() gosso.CookieOptions {
	return gosso.CookieOptions{Secure: func() bool { return config.GetAuth().OIDC.CookieSecure }}
}

type Auth struct {
	core     *gosso.Auth
	sessions SessionStore
	stop     context.CancelFunc
	wg       sync.WaitGroup
}

func New(sessions SessionStore) *Auth {
	auth := &Auth{sessions: sessions}
	auth.core = gosso.New(settings, NewDBSessions(sessions), callbackPath, "")
	auth.startCleanup()

	gosso.SetLogger(slog.Default())

	return auth
}

func (a *Auth) Enabled() bool {
	return a.core.Enabled()
}

func (a *Auth) StartLogin(c *echo.Context) error {
	return a.core.StartLogin(c.Response(), c.Request(), c.QueryParam("redirect"))
}

func (a *Auth) HandleCallback(c *echo.Context) (*gosso.LoginResult, *gosso.LoginFailure) {
	return a.core.HandleCallback(c.Response(), c.Request())
}

func (a *Auth) Authenticate(c *echo.Context) (*User, error) {
	return a.core.Authenticate(c.Response(), c.Request())
}

func (a *Auth) ClearSession(c *echo.Context) {
	a.core.ClearSession(c.Response(), c.Request())
}

func (a *Auth) SessionTTL() time.Duration {
	return a.core.SessionTTL()
}

func (a *Auth) Middleware() echo.MiddlewareFunc {
	return echo.WrapMiddleware(a.core.Middleware())
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

func settings() gosso.Config {
	oidc := config.GetAuth().OIDC

	return gosso.Config{
		Enabled:      oidc.Enabled,
		IssuerURL:    oidc.IssuerURL,
		ClientID:     oidc.ClientID,
		ClientSecret: oidc.ClientSecret,
		SessionTTL:   oidc.SessionTTL,
		CookieSecure: oidc.CookieSecure,
		AppURL:       config.GetAppUrl(),
	}
}
