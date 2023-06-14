package controller

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr"
	"go.uber.org/zap"
)

func (c *Controller) SendNotification(title string, msg string) {
	if c.env.NotificationURL == "" {
		return
	}
	url := fmt.Sprintf("%s&%s&Title=%s", c.env.NotificationURL, "parseMode=html", url.PathEscape(title))
	err := shoutrrr.Send(url, msg)
	if err != nil {
		zap.L().Error("cannot send notification", zap.String("msg", msg), zap.Error(err))
	}
}
