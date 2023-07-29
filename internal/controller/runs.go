package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Schemes
// @Tags		jobs
// @Produce	json
// @Success	200	{array}	database.Run
// @Router		/jobs/runs [get]
func (c *Controller) GetRuns(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetRuns())
}
