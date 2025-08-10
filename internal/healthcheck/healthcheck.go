package healthcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"gitlab.unjx.de/flohoss/gocron/config"
)

func SendStart() {
	if err := sendHttpRequest(config.GetHealthcheck().Start); err != nil {
		slog.Error("Failed to send healthcheck start event", "err", err.Error())
	}
}

func SendEnd() {
	if err := sendHttpRequest(config.GetHealthcheck().End); err != nil {
		slog.Error("Failed to send healthcheck end event", "err", err.Error())
	}
}

func SendFailure() {
	if err := sendHttpRequest(config.GetHealthcheck().Failure); err != nil {
		slog.Error("Failed to send healthcheck failure event", "err", err.Error())
	}
}

func sendHttpRequest(u config.Url) error {
	if u.Url == "" {
		return nil
	}
	parsedUrl, err := url.Parse(u.Url)
	if err != nil {
		return err
	}
	query := parsedUrl.Query()
	for key, value := range u.Params {
		query.Set(key, value)
	}
	parsedUrl.RawQuery = query.Encode()

	var bodyReader io.Reader
	if u.Body != "" {
		var jsonData any
		if err := json.Unmarshal([]byte(u.Body), &jsonData); err != nil {
			return fmt.Errorf("invalid JSON body: %w", err)
		}
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(config.GetHealthcheck().Type, parsedUrl.String(), bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.GetHealthcheck().Authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
