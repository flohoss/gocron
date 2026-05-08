package healthcheck

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flohoss/gocron/config"
	"github.com/spf13/viper"
)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := config.HealthCheck{
		Authorization: "Bearer abc",
		Type:          http.MethodPost,
	}
	loadHealthcheckConfig(t, hc)

	err := sendHttpRequest(config.Url{
		Url: server.URL,
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	loadHealthcheckConfig(t, config.HealthCheck{Type: http.MethodPost})

	err := sendHttpRequest(config.Url{Url: server.URL})
	if err == nil {
		t.Fatal("expected non-200 error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code: 500") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendStartEndFailure_CallConfiguredEndpoints(t *testing.T) {
	hits := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits[r.URL.Path]++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	loadHealthcheckConfig(t, config.HealthCheck{
		Type: http.MethodPost,
		Start: config.Url{Url: server.URL + "/start"},
		End: config.Url{Url: server.URL + "/end"},
		Failure: config.Url{Url: server.URL + "/failure"},
	})

	SendStart()
	SendEnd()
	SendFailure()

	if hits["/start"] != 1 || hits["/end"] != 1 || hits["/failure"] != 1 {
		t.Fatalf("unexpected endpoint hit counts: %#v", hits)
	}
}
