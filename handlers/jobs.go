package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/services"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
	"gitlab.unjx.de/flohoss/gobackup/views"
)

type JobService interface {
	GetQueries() *jobs.Queries
	ExecuteJobs()
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
	jobsAndRuns, _ := jh.JobService.GetQueries().ListJobsAndLatestRun(context.Background())
	templateJob := &services.TemplateJob{
		Name: "Home",
	}
	return renderView(c, views.HomeIndex(templateJob, views.Home(jobsAndRuns)))
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
	}

	templateJob.Commands, _ = jh.JobService.GetQueries().ListCommandsByJobID(context.Background(), job.ID)
	templateJob.Envs, _ = jh.JobService.GetQueries().ListEnvsByJobID(context.Background(), job.ID)

	runs, _ := jh.JobService.GetQueries().ListRunsByJobID(context.Background(), job.ID)
	for _, run := range runs {
		logs, _ := jh.JobService.GetQueries().ListLogsByRunID(context.Background(), run.ID)
		templateJob.Runs = append(templateJob.Runs, services.TemplateRun{
			StatusID:  run.StatusID,
			StartTime: run.StartTime,
			EndTime:   run.EndTime,
			Logs:      logs,
		})
	}

	return renderView(c, views.JobIndex(&templateJob, views.Job(&templateJob)))
}

func (jh *JobHandler) executeJobsHandler(c echo.Context) error {
	go jh.JobService.ExecuteJobs()
	return c.NoContent(http.StatusOK)
}
