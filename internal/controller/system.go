package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	system.Data
//	@Router		/system [get]
func (c *Controller) GetSystem(ctx echo.Context) error {
	stats := c.service.GetJobStats()
	return ctx.JSON(http.StatusOK, stats)
}

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{array}	database.SystemLog
//	@Router		/system/logs [get]
func (c *Controller) GetSystemLogs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetSystemLogs())
}
