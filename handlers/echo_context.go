package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
)

type echoContextKey struct{}

func echoContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), echoContextKey{}, c)))
		return next(c)
	}
}

func echoContextFrom(ctx context.Context) *echo.Context {
	c, _ := ctx.Value(echoContextKey{}).(*echo.Context)
	return c
}

func redirect(c *echo.Context, location string) error {
	return c.Redirect(http.StatusFound, location)
}
