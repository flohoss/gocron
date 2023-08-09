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

//	@Schemes
//	@Tags		jobs
//	@Accept		json
//	@Param		id	path	int	true	"Job ID"
//	@Success	200
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/jobs/{id}/runs [delete]
func (c *Controller) DeleteJobRuns(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	job := c.service.GetJob(id)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}
	c.service.DeleteJobRuns(job.ID)
	return ctx.NoContent(http.StatusOK)
}
