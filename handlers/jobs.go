package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/views"
)

type JobService interface {
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

func (jh *JobHandler) homeHandler(c echo.Context) error {
	return renderView(c, views.HomeIndex(views.Home(jh.Config)))
}
