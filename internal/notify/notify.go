package notify

import (
	"net/http"
	"strings"

	"github.com/labstack/gommon/log"
)

type Notifier struct {
	URL                  string
	Topic                string
	Token                string
	SendMessageOnSuccess bool
}

func New(url, topic, token string, sendMessageOnSuccess bool) *Notifier {
	return &Notifier{
		URL:                  url,
		Topic:                topic,
		Token:                token,
		SendMessageOnSuccess: sendMessageOnSuccess,
	}
}

type priorityLevel uint8

const (
	MIN priorityLevel = iota + 1
	LOW
	DEFAULT
	HIGH
	URGENT
)

func (p priorityLevel) String() string {
	switch p {
	case MIN:
		return "min"
	case LOW:
		return "low"
	case HIGH:
		return "high"
	case URGENT:
		return "urgent"
	default:
		return "default"
	}
}

func (n *Notifier) Send(title, message string, priority priorityLevel, tags []string) {
	if n.URL == "" || n.Topic == "" {
		return
	}
	req, _ := http.NewRequest("POST", n.URL+n.Topic, strings.NewReader(message))
	req.Header.Set("Title", title)
	req.Header.Set("Priority", priority.String())
	req.Header.Set("Tags", strings.Join(tags, ","))
	if n.Token != "" {
		req.Header.Set("Authorization", "Bearer "+n.Token)
	}

	body, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Failed to send notification (url: %s, topic: %s): %v", n.URL, n.Topic, err)
		return
	}
	defer body.Body.Close()
	if body.StatusCode != 200 {
		log.Warnf("Failed to send notification (url: %s, topic: %s): %s", n.URL, n.Topic, body.Status)
		return
	}
	log.Debugf("Notification sent (url: %s, topic: %s)", n.URL, n.Topic)
}
