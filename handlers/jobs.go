package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services"
	"github.com/flohoss/gocron/services/jobs"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

type JobService interface {
	GetQueries() *jobs.Queries
	GetParser() *cron.Parser
	GetHandler() echo.HandlerFunc
	IsIdle() bool
	ExecuteJobs(jobs []config.Job)
	ExecuteJob(job *config.Job)
	ListJobs() []services.JobView
	ListRuns(name string, limit int64) ([]services.RunView, error)
}

func NewJobHandler(js JobService) *JobHandler {
	return &JobHandler{
		JobService: js,
	}
}

type JobHandler struct {
	JobService JobService
}

func (jh *JobHandler) listJobsOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-jobs",
		Method:      http.MethodGet,
		Path:        "/api/jobs",
		Summary:     "Get jobs",
		Description: "Get jobs with run details but no logs.",
		Tags:        []string{"Jobs"},
	}
}

type Jobs struct {
	Body []services.JobView
}

func (jh *JobHandler) listJobsHandler(ctx context.Context, input *struct{}) (*Jobs, error) {
	jobs := jh.JobService.ListJobs()
	return &Jobs{Body: jobs}, nil
}

func (jh *JobHandler) listRunsOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-runs",
		Method:      http.MethodGet,
		Path:        "/api/runs/{job_name}",
		Summary:     "Get runs",
		Description: "Get runs with logs for a job.",
		Tags:        []string{"Runs"},
	}
}

type Runs struct {
	Body []services.RunView
}

func (jh *JobHandler) listRunsHandler(ctx context.Context, input *struct {
	Name  string `path:"job_name" minLength:"1" maxLength:"255" doc:"job name"`
	Limit int64  `query:"limit" default:"5" doc:"number of runs to return"`
}) (*Runs, error) {
	jobView, err := jh.JobService.ListRuns(input.Name, input.Limit)
	if err != nil {
		return nil, huma.Error404NotFound("Job not found")
	}
	return &Runs{Body: jobView}, nil
}

func (jh *JobHandler) executeJobsOperation() huma.Operation {
	return huma.Operation{
		OperationID: "post-jobs",
		Method:      http.MethodPost,
		Path:        "/api/jobs",
		Summary:     "Start jobs",
		Description: "Start all jobs in order of name.",
		Tags:        []string{"Jobs"},
	}
}

func (jh *JobHandler) executeJobsHandler(ctx context.Context, input *struct{}) (*struct{}, error) {
	go jh.JobService.ExecuteJobs([]config.Job{})
	return nil, nil
}

func (jh *JobHandler) executeJobOperation() huma.Operation {
	return huma.Operation{
		OperationID: "post-job",
		Method:      http.MethodPost,
		Path:        "/api/jobs/{name}",
		Summary:     "Start job",
		Description: "Start single job.",
		Tags:        []string{"Jobs"},
	}
}

func (jh *JobHandler) executeJobHandler(ctx context.Context, input *struct {
	Name string `path:"name" maxLength:"255" doc:"job name"`
}) (*struct{}, error) {
	job := config.GetJobByName(input.Name)
	if job == nil {
		return nil, huma.Error404NotFound("Job not found")
	}
	go jh.JobService.ExecuteJob(job)
	return nil, nil
}

func (jh *JobHandler) changeJobOperation() huma.Operation {
	return huma.Operation{
		OperationID: "put-job",
		Method:      http.MethodPut,
		Path:        "/api/jobs",
		Summary:     "Change job",
		Description: "Change single job.",
		Tags:        []string{"Jobs"},
	}
}

type Action string

const (
	Toggle             Action = "toggle"
	EnableAll          Action = "enable_all"
	DisableAll         Action = "disable_all"
	EnableScheduled    Action = "enable_scheduled"
	EnableNonScheduled Action = "enable_non_scheduled"
)

func (jh *JobHandler) changeJobHandler(ctx context.Context, input *struct {
	Name   string `query:"name" maxLength:"255" doc:"job name"`
	Action Action `query:"action" enum:"toggle,disable_all,enable_all,enable_scheduled,enable_non_scheduled" doc:"action to perform"`
}) (*struct{}, error) {
	switch input.Action {
	case DisableAll:
		config.DisableAllJobs()
	case EnableAll:
		config.EnableAllJobs()
	case EnableScheduled:
		config.EnableScheduledJobs()
	case EnableNonScheduled:
		config.EnableNonScheduledJobs()
	case Toggle:
		if err := config.ToggleDisabledJob(input.Name); err != nil {
			return nil, huma.Error404NotFound("Job not found")
		}
	}
	return nil, nil
}
