package notify

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
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
		return
	}
	resp.Body.Close()

}

type Notify struct {
	sender *router.ServiceRouter
}

func NewNotificationService(shoutrrrUrl string) *Notify {
	s, _ := shoutrrr.CreateSender(shoutrrrUrl)
	n := Notify{
		sender: s,
	}
	return &n
}

func (n *Notify) SendNotification(title string, msg string) {
	if n.sender == nil {
		return
	}
	err := n.sender.Send(msg, &types.Params{"title": title})
	if err != nil {
		slog.Error("cannot send notification", "msg", msg, "err", err)
	}
	slog.Debug("notification send")
}
