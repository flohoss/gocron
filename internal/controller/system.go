package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/internal/system"
)

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	system.SystemConfig
//	@Router		/system [get]
func (c *Controller) GetSystem(ctx echo.Context) error {
	config := system.SystemConf
	config.Config = *c.env
	return ctx.JSON(http.StatusOK, config)
}

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	database.JobStats
//	@Router		/system/stats [get]
func (c *Controller) GetSystemStats(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetJobStats())
}

//	@Schemes
//	@Tags		system
//	@Produce	json
//	@Success	200	{array}	database.SystemLog
//	@Router		/system/logs [get]
func (c *Controller) GetSystemLogs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetSystemLogs())
}
