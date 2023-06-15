package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/internal/models"
)

type RemotesData struct {
	Title   string
	Remotes []models.Remote
}

type FormRemoteData struct {
	Title              string
	Remote             *models.Remote
	RepositoryOptions  []models.SelectOption
	CompressionOptions []models.SelectOption
	Errors             map[string][]string
}

func (c *Controller) RenderRemotes(ctx echo.Context) error {
	var remotes []models.Remote
	c.orm.Find(&remotes)
	return ctx.Render(http.StatusOK, "remotes", RemotesData{Title: c.env.Identifier + " - Remotes", Remotes: remotes})
}

func (c *Controller) CreateRemoteConfiguration(ctx echo.Context) error {
	remote := new(models.Remote)
	err := c.bindToRequest(ctx, remote)
	if err == nil {
		data := FormRemoteData{Title: c.env.Identifier + " - New Remote", Remote: remote, RepositoryOptions: models.RepositoryTypes(), CompressionOptions: models.CompressionTypes()}
		if remote.ID != 0 {
			data.Title = c.env.Identifier + " - Edit Remote"
		}

		if errors := c.createOrUpdate(ctx, remote); len(errors) != 0 {
			data.Errors = errors
			return ctx.Render(http.StatusOK, "remotesForm", data)
		}
	}
	return c.RenderRemotes(ctx)
}

func (c *Controller) RenderRemoteForm(ctx echo.Context) error {
	remote := new(models.Remote)
	id := ctx.QueryParam("id")
	title := "New Remote"
	if id != "" {
		c.orm.Limit(1).Find(&remote, id)
		if remote.ID != 0 {
			title = "Edit Remote"
		}
	}
	return ctx.Render(http.StatusOK, "remotesForm", FormRemoteData{Title: c.env.Identifier + " - " + title, Remote: remote, RepositoryOptions: models.RepositoryTypes(), CompressionOptions: models.CompressionTypes()})
}

func (c *Controller) DeleteRemoteConfiguration(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.orm.Delete(&models.Remote{}, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (c *Controller) RemoteOptions() []models.SelectOption {
	var remotes []models.Remote
	c.orm.Find(&remotes)
	var temp []models.SelectOption
	for _, remote := range remotes {
		temp = append(temp, models.SelectOption{
			Name:  remote.Description,
			Value: int(remote.ID),
		})
	}
	return temp
}
