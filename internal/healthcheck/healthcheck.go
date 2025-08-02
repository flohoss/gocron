package healthcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.unjx.de/flohoss/gocron/config"
)

func SendStart() {
	sendHttpRequest(config.GetHealthcheck().Start)
}

func SendEnd() {
	sendHttpRequest(config.GetHealthcheck().End)
}

func SendFailure() {
	sendHttpRequest(config.GetHealthcheck().Failure)
}

func sendHttpRequest(u config.Url) error {
	parsedUrl, err := url.Parse(u.Url)
	if err != nil {
		return err
	}
	query := parsedUrl.Query()
	for key, value := range u.Params {
		query.Set(key, value)
	}
	parsedUrl.RawQuery = query.Encode()

	var jsonData interface{}
	if err := json.Unmarshal([]byte(u.Body), &jsonData); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(config.GetHealthcheck().Type, parsedUrl.String(), bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.GetHealthcheck().Authorization)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}
