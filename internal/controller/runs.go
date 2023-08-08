package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/database"
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
	start := time.Now().UnixMilli()
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	job := c.service.GetJob(id)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}
	c.service.DeleteJobRuns(job.ID)
	run := database.Run{
		JobID:     job.ID,
		StartTime: start,
		Status:    database.LogInfo,
	}
	c.service.CreateOrUpdate(&run)
	log := database.Log{
		RunID:       run.ID,
		LogType:     database.LogGeneral,
		LogSeverity: database.LogInfo,
		Message:     "runs cleared",
		CreatedAt:   time.Now().UnixMilli(),
	}
	c.service.CreateOrUpdate(&log)
	run.EndTime = time.Now().UnixMilli()
	c.service.CreateOrUpdate(&run)
	return ctx.NoContent(http.StatusOK)
}
