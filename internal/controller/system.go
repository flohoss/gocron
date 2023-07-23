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
	system.SystemData.Disk = system.DiskUsage()
	return ctx.JSON(http.StatusOK, system.SystemData)
}
