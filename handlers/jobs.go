package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

type JobService interface {
	GetQueries() *jobs.Queries
	GetParser() *cron.Parser
	GetHandler() echo.HandlerFunc
	IsIdle() bool
	ExecuteJobs(jobs []jobs.Job)
	ExecuteJob(job *jobs.Job)
	ListJobs(c echo.Context) error
	ListJob(c echo.Context) error
}

func NewJobHandler(js JobService, config *config.Config) *JobHandler {
	return &JobHandler{
		JobService: js,
	}
}

type JobHandler struct {
	JobService JobService
}

//	@Summary	List all jobs
//	@Produce	json
//	@Tags		jobs
//	@Success	200	{array}	jobs.JobsView
//	@Router		/jobs [get]
func (jh *JobHandler) listHandler(c echo.Context) error {
	return jh.JobService.ListJobs(c)
}

//	@Summary	List single job
//	@Produce	json
//	@Tags		jobs
//	@Param		name	path		string	true	"job id"
//	@Success	200		{object}	services.TemplateJob
//	@Failure	404		{object}	echo.HTTPError
//	@Router		/jobs/{name} [get]
func (jh *JobHandler) jobHandler(c echo.Context) error {
	return jh.JobService.ListJob(c)
}

//	@Summary	Run all jobs
//	@Produce	json
//	@Tags		jobs
//	@Success	200
//	@Router		/jobs [post]
func (jh *JobHandler) executeJobsHandler(c echo.Context) error {
	go jh.JobService.ExecuteJobs([]jobs.Job{})
	return c.NoContent(http.StatusOK)
}

//	@Summary	Run single job
//	@Produce	json
//	@Tags		jobs
//	@Param		name	path	string	true	"job id"
//	@Success	200
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/jobs/{name} [post]
func (jh *JobHandler) executeJobHandler(c echo.Context) error {
	name := c.Param("name")
	job, err := jh.JobService.GetQueries().GetJob(context.Background(), name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Job not found")
	}
	go jh.JobService.ExecuteJob(&job)
	return c.NoContent(http.StatusOK)
}
