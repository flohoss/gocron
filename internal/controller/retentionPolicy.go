package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//	@Schemes
//	@Tags		retention_policies
//	@Produce	json
//	@Success	200	{array}	database.SelectOption
//	@Router		/retention_policies [get]
func (c *Controller) GetRetentionPolicies(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetSelectOptionsRetentionPolicies())
}
