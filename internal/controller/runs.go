package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

//	@Schemes
//	@Tags		jobs
//	@Produce	json
//	@Param		id	path		int	true	"Job ID"
//	@Success	200	{array}		database.Run
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/jobs/{id}/runs [get]
func (c *Controller) GetRuns(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, c.service.GetRuns(id))
}
