package controller

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (c *Controller) sendHealthcheck(suffix string) {
	if c.env.HealthcheckURL == "" || c.env.HealthcheckUUID == "" {
		return
	}
	url := c.env.HealthcheckURL + c.env.HealthcheckUUID + suffix
	zap.L().Debug("Sending healthcheck", zap.String("url", url))
	var client = &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		zap.L().Error("cannot send healthcheck", zap.Error(err))
	}
	resp.Body.Close()
}
