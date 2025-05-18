package notify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/danielgtaylor/huma/v2"
	"github.com/enescakir/emoji"
	"github.com/labstack/gommon/log"
)

type Notifier struct {
	URL         string
	NotifyLevel log.Lvl
}

var (
	ErrNotificationDisabled = errors.New("notification is disabled, please set the NOTIFY_URL environment variable")
	ErrNotificationFailed   = errors.New("failed to send notification")
)

func New(url string, notifyLevel log.Lvl) *Notifier {
	return &Notifier{
		URL:         url,
		NotifyLevel: notifyLevel,
	}
}

func icon(level log.Lvl) emoji.Emoji {
	switch level {
	case log.DEBUG:
		return ""
	case log.INFO:
		return ""
	case log.WARN:
		return emoji.Warning
	case log.ERROR:
		return emoji.CrossMark
	default:
		return ""
	}
}

func (n *Notifier) Send(title, message string, level log.Lvl) error {
	// notification is disabled
	if n.URL == "" || level < n.NotifyLevel {
		return ErrNotificationDisabled
	}

	cmd := exec.Command("apprise", "-t", fmt.Sprintf("%s %s", icon(level), title), "-b", message, n.URL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w, output: %s, error: %v", ErrNotificationFailed, out, err)
	}
	return nil
}

func (n *Notifier) ExecuteNotifyOperation() huma.Operation {
	return huma.Operation{
		OperationID: "post-notify",
		Method:      http.MethodPost,
		Path:        "/api/notify",
		Summary:     "Test notification",
		Description: "Send a test notification to the configured URL.",
		Tags:        []string{"Notify"},
	}
}

func (n *Notifier) ExecuteNotifyHandler(ctx context.Context, input *struct{}) (*struct{}, error) {
	err := n.Send("Hello", "This is a test message from\nGoCron!", log.OFF)
	if err != nil {
		if errors.Is(err, ErrNotificationDisabled) {
			return nil, huma.Error412PreconditionFailed(err.Error())
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
