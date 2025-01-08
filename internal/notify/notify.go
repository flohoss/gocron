package notify

import (
	"net/http"
	"strings"

	"github.com/labstack/gommon/log"
)

type Notifier struct {
	URL   string
	Topic string
	Token string
}

func New(url, topic, token string) *Notifier {
	return &Notifier{
		URL:   url,
		Topic: topic,
		Token: token,
	}
}

func (n *Notifier) Send(title, message string, tags []string) {
	req, _ := http.NewRequest("POST", n.URL+n.Topic, strings.NewReader(message))
	req.Header.Set("Title", title)
	req.Header.Set("Priority", "urgent")
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
	log.Printf("Notification sent (url: %s, topic: %s)", n.URL, n.Topic)
}
