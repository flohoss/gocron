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
//	@Success	200	{array}	database.Job
//	@Router		/jobs [get]
func (c *Controller) GetJobs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetJobs())
}

//	@Schemes
//	@Tags		jobs
//	@Produce	json
//	@Param		id	path		int	true	"Job ID"
//	@Success	200	{object}	database.Job
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/jobs/{id} [get]
func (c *Controller) GetJob(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	job := c.service.GetJob(id)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}
	return ctx.JSON(http.StatusOK, job)
}

//	@Schemes
//	@Tags		jobs
//	@Accept		json
//	@Produce	json
//	@Param		job	body		database.Job	true	"Job"
//	@Success	200	{object}	database.Job
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Failure	500	{object}	echo.HTTPError
//	@Router		/jobs [put]
func (c *Controller) UpdateJob(ctx echo.Context) error {
	job := new(database.Job)
	if err := ctx.Bind(job); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	dbJob := c.service.GetJob(job.ID)
	if dbJob.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}
	jsonBlob, err := c.service.ValidateRequestBinding(ctx, job)
	if err != nil {
		return ctx.JSONBlob(http.StatusBadRequest, jsonBlob)
	}
	if err := c.service.CreateOrUpdateFromRequest(ctx, job); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, command := range dbJob.PreCommands {
		if !database.IsInArray(job.PreCommands, command) {
			c.service.DeleteCommand(command.ID)
		}
	}
	for _, command := range job.PreCommands {
		if err := c.service.CreateOrUpdate(&command); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	for _, command := range dbJob.PostCommands {
		if !database.IsInArray(job.PostCommands, command) {
			c.service.DeleteCommand(command.ID)
		}
	}
	for _, command := range job.PostCommands {
		if err := c.service.CreateOrUpdate(&command); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	job.Runs = dbJob.Runs
	return ctx.JSON(http.StatusOK, job)
}

//	@Schemes
//	@Tags		jobs
//	@Accept		json
//	@Produce	json
//	@Param		job	body		database.Job	true	"job"
//	@Success	201	{object}	database.Job
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	500	{object}	echo.HTTPError
//	@Router		/jobs [post]
func (c *Controller) CreateJob(ctx echo.Context) error {
	start := time.Now().UnixMilli()
	job := new(database.Job)
	if err := ctx.Bind(job); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	jsonBlob, err := c.service.ValidateRequestBinding(ctx, job)
	if err != nil {
		return ctx.JSONBlob(http.StatusBadRequest, jsonBlob)
	}
	if err := c.service.CreateOrUpdateFromRequest(ctx, job); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.service.CreateOrUpdate(&database.Run{
		JobID:     job.ID,
		StartTime: start,
		EndTime:   time.Now().UnixMilli(),
		Logs: []database.Log{{
			LogType:     database.LogGeneral,
			LogSeverity: database.LogInfo,
			Message:     "job created",
		}},
	})
	job.Status = database.LogInfo
	return ctx.JSON(http.StatusCreated, job)
}

//	@Schemes
//	@Tags		jobs
//	@Accept		json
//	@Param		id	path	int	true	"Job ID"
//	@Success	200
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/jobs/{id} [delete]
func (c *Controller) DeleteJob(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	job := c.service.GetJob(id)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}
	c.service.DeleteJob(job.ID)
	return ctx.NoContent(http.StatusOK)
}
