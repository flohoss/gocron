package handlers

import (
	"net/http"
	"slices"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/views"
)

type JobService interface {
	ExecuteJobs()
	ExecuteJob(id int)
}

func NewJobHandler(js JobService, config *config.Config) *JobHandler {
	return &JobHandler{
		JobService: js,
		Config:     config,
	}
}

type JobHandler struct {
	JobService JobService
	Config     *config.Config
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (jh *JobHandler) listHandler(c echo.Context) error {
	return renderView(c, views.HomeIndex(jh.Config, views.Home(jh.Config)))
}

func (jh *JobHandler) jobHandler(c echo.Context) error {
	name := c.Param("name")
	idx := slices.IndexFunc(jh.Config.Jobs, func(c config.Job) bool { return c.Name == name })

	if idx == -1 {
		return echo.ErrNotFound
	}
	job := &jh.Config.Jobs[idx]

	return renderView(c, views.JobIndex(jh.Config, views.Job(job)))
}

func (jh *JobHandler) executeJobsHandler(c echo.Context) error {
	jh.JobService.ExecuteJobs()
	return c.NoContent(http.StatusOK)
}
