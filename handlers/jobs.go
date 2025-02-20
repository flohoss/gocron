package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/internal/commands"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

type JobService interface {
	GetQueries() *jobs.Queries
	GetParser() *cron.Parser
	GetHandler() echo.HandlerFunc
	IsIdle() bool
	ExecuteJobs(jobs []jobs.Job)
	ExecuteJob(job *jobs.Job)
	ListJobs() []jobs.JobsView
	ListJob(name string) (*jobs.JobsView, error)
}

func NewJobHandler(js JobService, config *config.Config) *JobHandler {
	return &JobHandler{
		JobService: js,
	}
}

type JobHandler struct {
	JobService JobService
}

func (jh *JobHandler) listHandler(c echo.Context) error {
	jobs := jh.JobService.ListJobs()
	return c.JSON(http.StatusOK, jobs)
}

func (jh *JobHandler) jobHandler(c echo.Context) error {
	name := c.Param("name")
	jobView, err := jh.JobService.ListJob(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Job not found")
	}
	return c.JSON(http.StatusOK, jobView)
}

func (jh *JobHandler) executeJobsHandler(c echo.Context) error {
	go jh.JobService.ExecuteJobs([]jobs.Job{})
	return c.NoContent(http.StatusOK)
}

func (jh *JobHandler) executeJobHandler(c echo.Context) error {
	name := c.Param("name")
	job, err := jh.JobService.GetQueries().GetJob(context.Background(), name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Job not found")
	}
	go jh.JobService.ExecuteJob(&job)
	return c.NoContent(http.StatusOK)
}

func (jh *JobHandler) getVersions(ctx context.Context, input *struct{}) (*commands.Versions, error) {
	return commands.GetVersions(), nil
}
