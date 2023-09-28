package notify

import (
	"net/http"
	"strings"
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

type Notify struct {
	endpoint string
	token    string
	topic    string
}

func NewNotificationService(endpoint string, token string, topic string) *Notify {
	n := Notify{
		endpoint: endpoint,
		token:    token,
		topic:    topic,
	}
	return &n
}

func (n *Notify) SendNotification(title string, msg string) {
	if n.endpoint == "" || n.token == "" || n.topic == "" {
		return
	}
	req, _ := http.NewRequest("POST", n.endpoint+n.topic, strings.NewReader(msg))
	req.Header.Set("Title", title)
	http.DefaultClient.Do(req)
}
