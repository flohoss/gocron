package controller

import (
	"net/http"
	"strconv"

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
	if err := c.service.CreateOrUpdateJob(ctx, job); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, job)
}

//	@Schemes
//	@Tags		jobs
//	@Accept		json
//	@Produce	json
//	@Param		job	body		database.Job	true	"job"
//	@Success	200	{object}	database.Job
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	500	{object}	echo.HTTPError
//	@Router		/jobs [post]
func (c *Controller) CreateJob(ctx echo.Context) error {
	job := new(database.Job)
	if err := ctx.Bind(job); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.service.CreateOrUpdateJob(ctx, job); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, job)
}

//	@Schemes
//	@Tags		jobs
//	@Accept		json
//	@Produce	json
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

type CommandBody struct {
	Command string `json:"command" validate:"required"`
	JobID   uint64 `json:"job_id" validate:"omitempty,number"`
}

//	@Schemes
//	@Tags		commands
//	@Accept		json
//	@Produce	json
//	@Param		command	body	CommandBody	true	"Command body"
//	@Success	200
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/commands [post]
func (c *Controller) RunCommand(ctx echo.Context) error {
	cmdBody := new(CommandBody)
	if err := ctx.Bind(cmdBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(cmdBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	job := c.service.GetJob(cmdBody.JobID)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}

	switch cmdBody.Command {
	case "start":
		go c.runJob(func(job *database.Job) { c.runBackup(job) }, job)
	}

	return ctx.NoContent(http.StatusOK)
}
