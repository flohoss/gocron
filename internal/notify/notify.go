package notify

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func SendHealthcheck(url string, uuid string, suffix string) {
	if url == "" || uuid == "" {
		return
	}
	var client = &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Head(url + uuid + suffix)
	if err != nil {
		zap.L().Error("cannot send healthcheck", zap.Error(err))
	}
	resp.Body.Close()

}
