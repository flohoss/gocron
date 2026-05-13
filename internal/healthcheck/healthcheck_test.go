package healthcheck

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/flohoss/gocron/config"
	"github.com/spf13/viper"
)

// Use a custom transport so tests work in sandboxed CI where binding local ports is disallowed.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func withTestClient(t *testing.T, fn roundTripFunc) {
	t.Helper()
	origClient := httpClient
	httpClient = &http.Client{Transport: fn}
	t.Cleanup(func() {
		httpClient = origClient
	})
}

func loadHealthcheckConfig(t *testing.T, hc config.HealthCheck) {
	t.Helper()

	v := viper.New()
	v.Set("log_level", "info")
	v.Set("time_zone", "UTC")
	v.Set("delete_runs_after_days", 7)
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("terminal.allow_all_commands", true)
	v.Set("healthcheck", hc)

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
}

func TestSendHttpRequest_EmptyURLIsNoop(t *testing.T) {
	loadHealthcheckConfig(t, config.HealthCheck{Type: http.MethodPost})

	if err := sendHttpRequest(config.Url{}); err != nil {
		t.Fatalf("expected no error for empty url, got %v", err)
	}
}

func TestSendHttpRequest_SendsMethodHeadersQueryAndBody(t *testing.T) {
	type observedRequest struct {
		method        string
		authorization string
		query         map[string]string
		body          map[string]any
	}

	observed := observedRequest{}
	withTestClient(t, func(r *http.Request) (*http.Response, error) {
		observed.method = r.Method
		observed.authorization = r.Header.Get("Authorization")
		observed.query = map[string]string{
			"str":  r.URL.Query().Get("str"),
			"bool": r.URL.Query().Get("bool"),
			"int":  r.URL.Query().Get("int"),
			"flt":  r.URL.Query().Get("flt"),
		}

		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if err := json.Unmarshal(bodyBytes, &observed.body); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
		}, nil
	})

	hc := config.HealthCheck{
		Authorization: "Bearer abc",
		Type:          http.MethodPost,
	}
	loadHealthcheckConfig(t, hc)

	err := sendHttpRequest(config.Url{
		Url: "http://example.com/notify",
		Params: map[string]any{
			"str":  "value",
			"bool": true,
			"int":  42,
			"flt":  3.14,
		},
		Body: `{"name":"gocron"}`,
	})
	if err != nil {
		t.Fatalf("expected request to succeed, got %v", err)
	}

	if observed.method != http.MethodPost {
		t.Fatalf("unexpected method: got %q want %q", observed.method, http.MethodPost)
	}
	if observed.authorization != "Bearer abc" {
		t.Fatalf("unexpected authorization header: got %q", observed.authorization)
	}
	if observed.query["str"] != "value" || observed.query["bool"] != "true" || observed.query["int"] != "42" || observed.query["flt"] != "3.14" {
		t.Fatalf("unexpected query values: %#v", observed.query)
	}
	if observed.body["name"] != "gocron" {
		t.Fatalf("unexpected request body: %#v", observed.body)
	}
}

func TestSendHttpRequest_InvalidBodyReturnsError(t *testing.T) {
	loadHealthcheckConfig(t, config.HealthCheck{Type: http.MethodPost})

	err := sendHttpRequest(config.Url{
		Url:  "https://example.com",
		Body: "{not-json}",
	})
	if err == nil {
		t.Fatal("expected invalid json body error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid JSON body") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendHttpRequest_Non200ReturnsError(t *testing.T) {
	withTestClient(t, func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("boom")),
			Header:     make(http.Header),
		}, nil
	})

	loadHealthcheckConfig(t, config.HealthCheck{Type: http.MethodPost})

	err := sendHttpRequest(config.Url{Url: "http://example.com/fail"})
	if err == nil {
		t.Fatal("expected non-200 error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code: 500") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendStartEndFailure_CallConfiguredEndpoints(t *testing.T) {
	hits := map[string]int{}
	withTestClient(t, func(r *http.Request) (*http.Response, error) {
		hits[r.URL.Path]++
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
		}, nil
	})

	loadHealthcheckConfig(t, config.HealthCheck{
		Type: http.MethodPost,
		Start: config.Url{Url: "http://example.com/start"},
		End: config.Url{Url: "http://example.com/end"},
		Failure: config.Url{Url: "http://example.com/failure"},
	})

	SendStart()
	SendEnd()
	SendFailure()

	if hits["/start"] != 1 || hits["/end"] != 1 || hits["/failure"] != 1 {
		t.Fatalf("unexpected endpoint hit counts: %#v", hits)
	}
}
