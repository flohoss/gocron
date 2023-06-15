package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/internal/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"gitlab.unjx.de/flohoss/gobackup/internal/models"
	"gitlab.unjx.de/flohoss/gobackup/internal/validate"
	"gorm.io/gorm"
)

type Controller struct {
	orm           *gorm.DB
	env           *env.Config
	schedule      *cron.Cron
	SSE           *sse.Server
	Versions      Versions
	Configuration Configuration
	JobRunning    bool
}

func NewController(env *env.Config) *Controller {
	db := database.NewDatabaseConnection("sqlite.db")

	db.AutoMigrate(&models.Job{})
	db.AutoMigrate(&models.Remote{})
	db.AutoMigrate(&models.Log{})
	db.AutoMigrate(&models.SystemLog{})

	ctrl := Controller{orm: db, env: env, SSE: sse.New()}
	ctrl.setupSchedule()
	ctrl.setupEventChannel()
	ctrl.setVersions()
	return &ctrl
}

func (c *Controller) setupSchedule() {
	c.schedule = cron.New()
	if c.env.BackupCron != "" {
		c.schedule.AddFunc(c.env.BackupCron, func() {
			c.runBackups()
		})
	}
	if c.env.CleanupCron != "" {
		c.schedule.AddFunc(c.env.CleanupCron, func() {
			c.runPrunes()
		})
	}
	if c.env.CheckCron != "" {
		c.schedule.AddFunc(c.env.CheckCron, func() {
			c.runChecks()
		})
	}
	c.schedule.Start()
}

func (c *Controller) bindToRequest(ctx echo.Context, value interface{}) error {
	if err := ctx.Bind(value); err != nil {
		return err
	}
	return nil
}

func (c *Controller) createOrUpdate(ctx echo.Context, value interface{}) map[string][]string {
	tmp := make(map[string][]string)
	if err := ctx.Validate(value); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			tmp[err.Field()] = append(tmp[err.Field()], validate.MessageForError(err))
		}
		return tmp
	}
	if err := c.orm.Save(value).Error; err != nil {
		tmp["Global"] = append(tmp["Global"], "Could not save form. Please try again.")
		return tmp
	}
	return nil
}
