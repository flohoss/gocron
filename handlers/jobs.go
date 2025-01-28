package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/internal/commands"
	"gitlab.unjx.de/flohoss/gobackup/services"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
	"gitlab.unjx.de/flohoss/gobackup/views"
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

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (jh *JobHandler) listHandler(c echo.Context) error {
	templateJob := &services.TemplateJob{Name: "Home"}

	resultSet, _ := jh.JobService.GetQueries().GetJobsView(context.Background())
	jobsAmount := len(resultSet)
	for i := 0; i < jobsAmount; i++ {
		resultSet[i].Runs, _ = jh.JobService.GetQueries().GetRunsViewHome(context.Background(), resultSet[i].ID)
	}
	return renderView(c, views.HomeIndex(templateJob, commands.GetVersions(), jh.JobService.IsIdle(), views.Home(resultSet, jh.JobService.GetParser())))
}

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

	return renderView(c, views.JobIndex(&templateJob, jh.JobService.IsIdle(), views.Job(&templateJob)))
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
