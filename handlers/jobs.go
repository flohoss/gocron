package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gocron/internal/commands"
	"gitlab.unjx.de/flohoss/gocron/internal/notify"
	"gitlab.unjx.de/flohoss/gocron/services/jobs"
)

type JobService interface {
	GetQueries() *jobs.Queries
	GetParser() *cron.Parser
	GetHandler() echo.HandlerFunc
	IsIdle() bool
	ExecuteJobs(jobs []jobs.Job)
	ExecuteJob(job *jobs.Job)
	ListJobs() []jobs.JobsView
	ListJob(name string, limit int64) (*jobs.JobsView, error)
}

func NewJobHandler(js JobService, n *notify.Notifier) *JobHandler {
	return &JobHandler{
		JobService: js,
		Notify:     n,
	}
}

type JobHandler struct {
	JobService JobService
	Notify     *notify.Notifier
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
	Body []jobs.JobsView
}

func (jh *JobHandler) listJobsHandler(ctx context.Context, input *struct{}) (*Jobs, error) {
	jobs := jh.JobService.ListJobs()
	return &Jobs{Body: jobs}, nil
}

func (jh *JobHandler) listJobOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-job",
		Method:      http.MethodGet,
		Path:        "/api/jobs/{name}",
		Summary:     "Get job",
		Description: "Get job with run details and logs.",
		Tags:        []string{"Jobs"},
	}
}

type Job struct {
	Body *jobs.JobsView
}

func (jh *JobHandler) listJobHandler(ctx context.Context, input *struct {
	Name  string `path:"name" maxLength:"20" doc:"job name"`
	Limit int64  `query:"limit" default:"5" doc:"number of runs to return"`
}) (*Job, error) {
	jobView, err := jh.JobService.ListJob(input.Name, input.Limit)
	if err != nil {
		return nil, huma.Error404NotFound("Job not found")
	}
	return &Job{Body: jobView}, nil
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
	go jh.JobService.ExecuteJobs([]jobs.Job{})
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
	Name string `path:"name" maxLength:"20" doc:"job name"`
}) (*struct{}, error) {
	job, err := jh.JobService.GetQueries().GetJob(context.Background(), input.Name)
	if err != nil {
		return nil, huma.Error404NotFound("Job not found")
	}
	go jh.JobService.ExecuteJob(&job)
	return nil, nil
}

func (jh *JobHandler) getVersionsOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-versions",
		Method:      http.MethodGet,
		Path:        "/api/versions",
		Summary:     "Get installed versions",
		Description: "Get installed versions of software.",
		Tags:        []string{"Software"},
	}
}

type Versions struct {
	Body *commands.Versions
}

func (jh *JobHandler) getVersionsHandler(ctx context.Context, input *struct{}) (*Versions, error) {
	versions := commands.GetVersions()
	return &Versions{Body: versions}, nil
}
