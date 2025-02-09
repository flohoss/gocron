package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/services"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

type JobService interface {
	GetQueries() *jobs.Queries
	GetParser() *cron.Parser
	GetHandler() echo.HandlerFunc
	IsIdle() bool
	ExecuteJobs(jobs []jobs.Job)
	ExecuteJob(job *jobs.Job)
}

func NewJobHandler(js JobService, config *config.Config) *JobHandler {
	return &JobHandler{
		JobService: js,
	}
}

type JobHandler struct {
	JobService JobService
}

// @Summary	List all jobs
// @Produce	json
// @Tags		jobs
// @Success	200	{array}	[]jobs.JobsView
// @Router		/jobs [get]
func (jh *JobHandler) listHandler(c echo.Context) error {
	resultSet, _ := jh.JobService.GetQueries().GetJobsView(context.Background())
	jobsAmount := len(resultSet)
	for i := 0; i < jobsAmount; i++ {
		resultSet[i].Runs, _ = jh.JobService.GetQueries().GetRunsViewHome(context.Background(), resultSet[i].ID)
	}
	return c.JSON(http.StatusOK, resultSet)
}

// @Summary	List single job
// @Produce	json
// @Tags		jobs
// @Param name  path string true "job id"
// @Success	200	{object}	services.TemplateJob
// @Router		/jobs/{name} [get]
func (jh *JobHandler) jobHandler(c echo.Context) error {
	name := c.Param("name")

	job, err := jh.JobService.GetQueries().GetJob(context.Background(), name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	templateJob := services.TemplateJob{
		Name: job.Name,
		Cron: job.Cron,
		Job:  job,
	}

	templateJob.Commands, _ = jh.JobService.GetQueries().ListCommandsByJobID(context.Background(), job.ID)
	templateJob.Envs, _ = jh.JobService.GetQueries().ListEnvsByJobID(context.Background(), job.ID)

	runs, _ := jh.JobService.GetQueries().GetRunsViewDetail(context.Background(), job.ID)
	for _, run := range runs {
		logs, _ := jh.JobService.GetQueries().ListLogsByRunID(context.Background(), run.ID)
		run.Logs = logs
		templateJob.Runs = append(templateJob.Runs, run)
	}

	return c.JSON(http.StatusOK, templateJob)
}

func (jh *JobHandler) executeJobsHandler(c echo.Context) error {
	go jh.JobService.ExecuteJobs([]jobs.Job{})
	return c.NoContent(http.StatusOK)
}

func (jh *JobHandler) executeJobHandler(c echo.Context) error {
	name := c.Param("name")
	job, err := jh.JobService.GetQueries().GetJob(context.Background(), name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	go jh.JobService.ExecuteJob(&job)
	return c.NoContent(http.StatusOK)
}
