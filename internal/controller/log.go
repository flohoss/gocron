package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"gitlab.unjx.de/flohoss/gobackup/internal/models"
	"go.uber.org/zap"
)

type LogsData struct {
	Title      string
	JobOptions []models.SelectOption
}

func (c *Controller) addLogEntry(log models.Log, description string) {
	if err := c.orm.Save(&log).Error; err != nil {
		zap.L().Error("cannot save log entry", zap.String("job", log.Job.Description), zap.Error(err))
	}
	json, _ := json.Marshal(models.LogMessage{
		Description: description,
		Type:        log.Type,
		Topic:       log.Topic,
		Message:     log.Message,
		CreatedAt:   log.CreatedAt,
	})
	if log.Type != models.Info {
		lastLines := regexp.MustCompile(`(?:.*\n?){3}$`)
		splitted := lastLines.FindAllString(log.Message, -1)
		if splitted != nil {
			log.Message = splitted[len(splitted)-1]
		}
		c.SendNotification(fmt.Sprintf("%s %s %s", log.Type.Emoji().String(), c.env.Identifier, log.Type.String()), fmt.Sprintf("<b>Job:</b> <pre>%s</pre> <b>Topic:</b> <pre>%s</pre>\n<pre>%s</pre>", description, log.Topic, log.Message))
	}
	c.SSE.Publish(fmt.Sprintf("%s%d", EventLog, log.JobID), &sse.Event{Data: json})
	c.SSE.Publish(EventLog, &sse.Event{Data: json})
}

func (c *Controller) RenderLogs(ctx echo.Context) error {
	options := c.getJobOptions()
	return ctx.Render(http.StatusOK, "logs", LogsData{Title: c.env.Identifier + " - Logs", JobOptions: options})
}

func (c *Controller) GetLogs(ctx echo.Context) error {
	var logs []models.Log
	sysLogs := []models.SystemLog{}
	logMessages := []models.LogMessage{}

	limit := 20
	qLimit := ctx.QueryParam("limit")
	if qLimit != "" {
		var err error
		limit, err = strconv.Atoi(qLimit)
		if err != nil || limit < 0 {
			return ctx.NoContent(http.StatusBadRequest)
		} else if limit == 0 {
			limit = -1
		}
	}

	id := ctx.QueryParam("id")
	if id == "system" {
		c.orm.Order("created_at DESC").Limit(limit).Find(&sysLogs)

		for i := 0; i < len(sysLogs); i++ {
			logMessages = append(logMessages, models.LogMessage{
				Description: "System",
				Type:        sysLogs[i].Type,
				Topic:       sysLogs[i].Topic,
				Message:     sysLogs[i].Message,
				CreatedAt:   sysLogs[i].CreatedAt,
			})
		}
	} else {
		if id != "" {
			if _, err := strconv.Atoi(id); err != nil {
				return ctx.NoContent(http.StatusBadRequest)
			}
		}
		if id == "" {
			c.orm.Preload("Job").Order("created_at DESC").Limit(limit).Find(&logs)
		} else {
			c.orm.Where("job_id = ?", id).Preload("Job").Order("created_at DESC").Limit(limit).Find(&logs)
		}

		for i := 0; i < len(logs); i++ {
			logMessages = append(logMessages, models.LogMessage{
				Description: logs[i].Job.Description,
				Type:        logs[i].Type,
				Topic:       logs[i].Topic,
				Message:     logs[i].Message,
				CreatedAt:   logs[i].CreatedAt,
			})
		}
	}

	return ctx.JSON(http.StatusOK, logMessages)
}

func (c *Controller) addSystemLogEntry(log models.SystemLog) {
	if err := c.orm.Save(&log).Error; err != nil {
		zap.L().Error("cannot save log entry", zap.Error(err))
	}
	json, _ := json.Marshal(models.LogMessage{
		Description: "System",
		Type:        log.Type,
		Topic:       log.Topic,
		Message:     log.Message,
		CreatedAt:   log.CreatedAt,
	})
	c.SSE.Publish(EventSystem, &sse.Event{Data: json})
}
