package healthcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type HealthCheck struct {
	Authorization string `yaml:"authorization"`
	Type          string `validate:"omitempty,oneof=HEAD GET POST" yaml:"type"`
	Start         Url    `yaml:"start"`
	End           Url    `yaml:"end"`
	Failure       Url    `yaml:"failure"`
}

type Url struct {
	Url    string            `yaml:"url"`
	Params map[string]string `yaml:"params"`
	Body   string            `yaml:"body"`
}

func (h *HealthCheck) SendStart() {
	if h.Start.Url == "" {
		return
	}
	h.sendHttpRequest(&h.Start)
}

func (h *HealthCheck) SendEnd(body string) {
	if h.End.Url == "" {
		return
	}
	if h.End.Body == "" {
		h.End.Body = body
	}
	h.sendHttpRequest(&h.End)
}

func (h *HealthCheck) SendFailure() {
	if h.Failure.Url == "" {
		return
	}
	h.sendHttpRequest(&h.Failure)
}

func (h *HealthCheck) sendHttpRequest(u *Url) error {
	url, err := u.URLWithParams()
	if err != nil {
		return err
	}
	body, err := u.JSONBodyReader()
	if err != nil {
		return err
	}
	if h.Type == "" {
		h.Type = "POST"
	}
	req, err := http.NewRequest(h.Type, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", h.getAuthorization())
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

func (h *HealthCheck) getAuthorization() string {
	return os.ExpandEnv(h.Authorization)
}

func (u *Url) getUrl() (*url.URL, error) {
	return url.Parse(os.ExpandEnv(u.Url))
}

func (u *Url) JSONBodyReader() (io.Reader, error) {
	var jsonData interface{}

	if err := json.Unmarshal([]byte(u.Body), &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal JSON: %w", err)
	}

	return bytes.NewReader(jsonBytes), nil
}

func (u *Url) URLWithParams() (string, error) {
	parsedUrl, err := u.getUrl()
	if err != nil {
		return "", err
	}

	query := parsedUrl.Query()
	for key, value := range u.Params {
		query.Set(key, value)
	}
	parsedUrl.RawQuery = query.Encode()

	return parsedUrl.String(), nil
}
