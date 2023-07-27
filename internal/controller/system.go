package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/internal/system"
)

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	system.Data
//	@Router		/system [get]
func (c *Controller) GetSystem(ctx echo.Context) error {
	data := system.SystemData
	data.Stats = c.service.GetJobStats()
	return ctx.JSON(http.StatusOK, data)
}

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{array}	database.SystemLog
//	@Router		/system/logs [get]
func (c *Controller) GetSystemLogs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetSystemLogs())
}
