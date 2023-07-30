package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/containrrr/shoutrrr"
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

func SendNotification(shoutrrrUrl string, title string, msg string) {
	if shoutrrrUrl == "" {
		return
	}
	url := fmt.Sprintf("%s&%s&Title=%s", shoutrrrUrl, "parseMode=html", url.PathEscape(title))
	err := shoutrrr.Send(url, msg)
	if err != nil {
		zap.L().Error("cannot send notification", zap.String("msg", msg), zap.Error(err))
	}
}
