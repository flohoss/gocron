package notify

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/containrrr/shoutrrr"
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
		slog.Error("cannot send healthcheck", "err", err)
	}
	resp.Body.Close()

}

type Notify struct {
	shoutrrrUrl string
}

func NewNotificationService(shoutrrrUrl string) *Notify {
	n := Notify{
		shoutrrrUrl: shoutrrrUrl,
	}
	return &n
}

func (n *Notify) SendNotification(title string, msg string) {
	if n.shoutrrrUrl == "" {
		return
	}
	s := fmt.Sprintf("%s&parseMode=html&Title=%s", n.shoutrrrUrl, url.PathEscape(title))
	err := shoutrrr.Send(s, msg)
	if err != nil {
		slog.Error("cannot send notification", "msg", msg, "err", err)
	}
	slog.Debug("notification send")
}
