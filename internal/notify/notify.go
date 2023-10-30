package notify

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
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
		os.Exit(1)
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
	if n.endpoint == "" || n.topic == "" {
		return
	}
	req, _ := http.NewRequest("POST", n.endpoint+n.topic, strings.NewReader(msg))
	req.Header.Set("Title", title)
	if n.token != "" {
		req.Header.Set("Authorization", "Bearer "+n.token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return
	}
	slog.Info("notification send", "resp", string(b))
}
