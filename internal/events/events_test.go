package events

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
)

func newTestEvent(onSubscribe func(streamID string)) *Event {
	if onSubscribe == nil {
		onSubscribe = func(string) {}
	}
	return New(onSubscribe)
}

func newEventServer(t *testing.T, e *Event) *httptest.Server {
	t.Helper()

	router := echo.New()
	router.GET("/api/events", e.GetHandler())

	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	return srv
}

func (e *Event) subscriberCount(streamID string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return len(e.streams[streamID])
}

func TestNew_CreatesStreams(t *testing.T) {
	e := newTestEvent(nil)

	if !e.hasStream(EventStatus) {
		t.Fatalf("expected stream %q to exist", EventStatus)
	}
	if !e.hasStream(CommandEvent) {
		t.Fatalf("expected stream %q to exist", CommandEvent)
	}
}

func TestSendJobEvent_DoesNotPanic(t *testing.T) {
	e := newTestEvent(nil)
	// Publish to streams with no subscribers must not panic
	e.SendJobEvent(true, "run-1", []string{"job-a"})
	e.SendJobEvent(false, nil, nil)
}

func TestSendCommandEvent_DoesNotPanic(t *testing.T) {
	e := newTestEvent(nil)
	e.SendCommandEvent(2, "echo hello")
	e.SendCommandEvent(0, "")
}

func TestGetHandler_ReturnsNonNilHandler(t *testing.T) {
	e := newTestEvent(nil)
	if e.GetHandler() == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestOnSubscribe_CalledForKnownStream(t *testing.T) {
	called := make(chan string, 1)
	e := newTestEvent(func(streamID string) {
		select {
		case called <- streamID:
		default:
		}
	})

	srv := newEventServer(t, e)

	ctx := t.Context()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/events?stream="+EventStatus, nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			buf := make([]byte, 64)
			_, _ = resp.Body.Read(buf)
		}
	}()

	select {
	case streamID := <-called:
		if streamID != EventStatus {
			t.Fatalf("expected stream %q, got %q", EventStatus, streamID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("onSubscribe was not called for a known stream")
	}
}

func TestGetHandler_UnknownStreamReturnsNotFound(t *testing.T) {
	e := newTestEvent(nil)

	srv := newEventServer(t, e)

	resp, err := http.Get(srv.URL + "/api/events?stream=unknown")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestSendJobEvent_DeliversToSubscriber(t *testing.T) {
	e := newTestEvent(nil)
	subscriber := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, subscriber)

	e.SendJobEvent(true, nil, nil)

	select {
	case data := <-subscriber:
		if !strings.Contains(string(data), `"idle":true`) {
			t.Fatalf("unexpected payload: %s", data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("subscriber did not receive the event")
	}
}

func TestSendJobEvent_SlowSubscriberDoesNotBlockOthers(t *testing.T) {
	e := newTestEvent(nil)

	slow := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, slow)
	healthy := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, healthy)

	for i := 0; i < subscriberBuffer*4; i++ {
		e.SendJobEvent(false, i, nil)
	}

	drain(healthy)

	e.SendJobEvent(true, "marker", nil)

	select {
	case data := <-healthy:
		if !strings.Contains(string(data), "marker") {
			t.Fatalf("expected marker, got %s", data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("healthy subscriber never received the marker while a slow subscriber was stalled")
	}
}

func drain(subscriber chan []byte) {
	for {
		select {
		case <-subscriber:
		default:
			return
		}
	}
}

func TestUnsubscribe_ClosesChannel(t *testing.T) {
	e := newTestEvent(nil)
	subscriber := e.subscribe(EventStatus)

	e.unsubscribe(EventStatus, subscriber)

	if _, open := <-subscriber; open {
		t.Fatal("expected subscriber channel to be closed")
	}
}

func TestUnsubscribe_UnknownSubscriberIsNoop(t *testing.T) {
	e := newTestEvent(nil)
	subscriber := e.subscribe(EventStatus)
	e.unsubscribe(EventStatus, subscriber)

	// Second unsubscribe must not panic on a closed channel.
	e.unsubscribe(EventStatus, subscriber)
}

func newTestEventWithKeepAlive(keepAlive time.Duration) *Event {
	e := newTestEvent(nil)
	e.keepAlive = keepAlive
	return e
}

// startStream opens a streaming request and returns the response plus a reader.
func startStream(t *testing.T, srv *httptest.Server, streamID string) (context.CancelFunc, *http.Response, *bufio.Reader) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/events?stream="+streamID, nil)
	if err != nil {
		cancel()
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		t.Fatal(err)
	}

	return cancel, resp, bufio.NewReader(resp.Body)
}

func readUntil(t *testing.T, reader *bufio.Reader, match string, timeout time.Duration) string {
	t.Helper()

	deadline := time.Now().Add(timeout)
	var seen []string
	for time.Now().Before(deadline) {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		seen = append(seen, line)
		if strings.Contains(line, match) {
			return strings.Join(seen, "")
		}
	}

	t.Fatalf("did not observe %q within %s; saw %q", match, timeout, strings.Join(seen, ""))
	return ""
}

// readFrame reads lines until a complete SSE frame (data line plus terminating
// blank line) is consumed and returns the frame as received. Comment-only
// keep-alive frames are skipped.
func readFrame(t *testing.T, reader *bufio.Reader, timeout time.Duration) string {
	t.Helper()

	type result struct {
		data string
		err  error
	}
	done := make(chan result, 1)

	go func() {
		var frame strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				done <- result{"", err}
				return
			}
			if strings.HasPrefix(line, ":") {
				continue
			}
			frame.WriteString(line)
			if line == "\n" {
				done <- result{frame.String(), nil}
				return
			}
		}
	}()

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("failed to read frame: %v", r.err)
		}
		return r.data
	case <-time.After(timeout):
		t.Fatal("timed out waiting for an SSE frame")
		return ""
	}
}

func TestGetHandler_SetsEventStreamHeaders(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	cancel, resp, _ := startStream(t, srv, EventStatus)
	defer cancel()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("expected text/event-stream, got %q", got)
	}
	if got := resp.Header.Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("expected no-cache, got %q", got)
	}
	if got := resp.Header.Get("Connection"); got != "keep-alive" {
		t.Fatalf("expected keep-alive, got %q", got)
	}
	// nginx buffers proxied responses by default; this header makes it stream
	// the event stream instead. Other proxies ignore the response header.
	if got := resp.Header.Get("X-Accel-Buffering"); got != "no" {
		t.Fatalf("expected X-Accel-Buffering no, got %q", got)
	}
}

func TestGetHandler_WritesDataFrameAndFlushes(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	cancel, resp, reader := startStream(t, srv, EventStatus)
	defer cancel()
	defer resp.Body.Close()

	e.SendJobEvent(true, "run-42", nil)

	frame := readFrame(t, reader, 2*time.Second)
	if !strings.HasPrefix(frame, "data: ") {
		t.Fatalf("expected data frame, got %q", frame)
	}
	if !strings.HasSuffix(frame, "\n\n") {
		t.Fatalf("expected event terminating blank line, got %q", frame)
	}
}

func TestGetHandler_KeepAliveComment(t *testing.T) {
	e := newTestEventWithKeepAlive(20 * time.Millisecond)
	srv := newEventServer(t, e)

	cancel, resp, reader := startStream(t, srv, EventStatus)
	defer cancel()
	defer resp.Body.Close()

	readUntil(t, reader, ": keep-alive", 2*time.Second)
}

func TestGetHandler_DisconnectRemovesSubscriber(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	cancel, resp, _ := startStream(t, srv, EventStatus)

	if got := e.subscriberCount(EventStatus); got != 1 {
		t.Fatalf("expected 1 subscriber, got %d", got)
	}

	cancel()
	resp.Body.Close()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if e.subscriberCount(EventStatus) == 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected subscriber to be removed after disconnect, got %d", e.subscriberCount(EventStatus))
}

func TestWriteFrame_TimesOutWhenClientStalled(t *testing.T) {
	// A stalled client that never reads must not hold the handler goroutine
	// forever: once the write deadline passes the frame write fails and the
	// handler unsubscribes.
	w := &blockingResponseWriter{}
	controller := http.NewResponseController(w)

	start := time.Now()
	err := writeFrame(w, controller, "data: x\n\n", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected a write error after the deadline passed")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("writeFrame blocked far past its deadline: %s", elapsed)
	}
}

type blockingResponseWriter struct {
	deadline time.Time
}

func (w *blockingResponseWriter) Header() http.Header { return http.Header{} }

func (w *blockingResponseWriter) WriteHeader(int) {}

func (w *blockingResponseWriter) Write(p []byte) (int, error) {
	// Simulate a client that never reads: block until the write deadline.
	timer := time.NewTimer(time.Until(w.deadline))
	defer timer.Stop()
	<-timer.C
	return 0, fmt.Errorf("write timeout after %d bytes", len(p))
}

func (w *blockingResponseWriter) Flush() {}

func (w *blockingResponseWriter) SetWriteDeadline(deadline time.Time) error {
	w.deadline = deadline
	return nil
}

func TestGetHandler_MultipleSubscribersAllReceive(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	const subscribers = 5
	readers := make([]*bufio.Reader, 0, subscribers)
	for i := 0; i < subscribers; i++ {
		cancel, resp, reader := startStream(t, srv, EventStatus)
		defer cancel()
		defer resp.Body.Close()
		readers = append(readers, reader)
	}

	deadline := time.Now().Add(2 * time.Second)
	for e.subscriberCount(EventStatus) != subscribers && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}

	e.SendJobEvent(true, "broadcast", nil)

	for i, reader := range readers {
		readUntil(t, reader, `"run":"broadcast"`, 2*time.Second)
		_ = i
	}
}

func TestConcurrentPublishAndSubscribe(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	const (
		publishers   = 8
		perPublisher = 200
		connections  = 10
	)

	var publishersWG sync.WaitGroup
	for p := 0; p < publishers; p++ {
		publishersWG.Add(1)
		go func(id int) {
			defer publishersWG.Done()
			for i := 0; i < perPublisher; i++ {
				e.SendJobEvent(i%2 == 0, id, nil)
				e.SendCommandEvent(i%4, "echo")
			}
		}(p)
	}

	var connectionsWG sync.WaitGroup
	for c := 0; c < connections; c++ {
		connectionsWG.Add(1)
		go func() {
			defer connectionsWG.Done()

			for range 20 {
				ctx, cancel := context.WithCancel(context.Background())
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/events?stream="+EventStatus, nil)
				if err != nil {
					cancel()
					continue
				}
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					cancel()
					continue
				}

				deadline := time.Now().Add(200 * time.Millisecond)
				for e.subscriberCount(EventStatus) == 0 && time.Now().Before(deadline) {
					time.Sleep(time.Millisecond)
				}

				cancel()
				resp.Body.Close()
			}
		}()
	}

	publishersWG.Wait()
	connectionsWG.Wait()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if e.subscriberCount(EventStatus) == 0 && e.subscriberCount(CommandEvent) == 0 {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("subscriber leak: status=%d command=%d", e.subscriberCount(EventStatus), e.subscriberCount(CommandEvent))
}

func TestPublishNeverBlocksWithStalledSubscriber(t *testing.T) {
	e := newTestEvent(nil)

	stalled := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, stalled)

	published := make(chan struct{})
	go func() {
		for i := 0; i < subscriberBuffer*10; i++ {
			e.SendJobEvent(false, i, nil)
		}
		close(published)
	}()

	select {
	case <-published:
	case <-time.After(3 * time.Second):
		t.Fatal("SendJobEvent blocked while a subscriber was stalled")
	}

	healthy := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, healthy)

	e.SendJobEvent(true, "fresh", nil)

	select {
	case data := <-healthy:
		if !strings.Contains(string(data), "fresh") {
			t.Fatalf("expected fresh event, got %s", data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("fresh subscriber did not receive an event")
	}
}

func TestPublish_DropsOldestKeepingNewestForStalledSubscriber(t *testing.T) {
	e := newTestEvent(nil)
	subscriber := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, subscriber)

	for i := 0; i < subscriberBuffer*3; i++ {
		e.SendJobEvent(false, i, nil)
	}
	e.SendJobEvent(true, "latest", nil)

	var last []byte
	for {
		select {
		case data := <-subscriber:
			last = data
			continue
		default:
		}
		break
	}

	if last == nil {
		t.Fatal("expected the subscriber buffer to retain events")
	}
	if !strings.Contains(string(last), "latest") {
		t.Fatalf("expected the newest event to be retained, got %s", last)
	}
}

func TestSendCommandEvent_OnlyToCommandStream(t *testing.T) {
	e := newTestEvent(nil)
	status := e.subscribe(EventStatus)
	defer e.unsubscribe(EventStatus, status)
	command := e.subscribe(CommandEvent)
	defer e.unsubscribe(CommandEvent, command)

	e.SendCommandEvent(3, "ls -la")

	select {
	case data := <-command:
		if !strings.Contains(string(data), "ls -la") {
			t.Fatalf("unexpected command payload: %s", data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("command subscriber did not receive the event")
	}

	select {
	case data := <-status:
		t.Fatalf("status subscriber should not receive command events, got %s", data)
	default:
	}
}

func TestGetHandler_ServesThroughEchoRouter(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	cancel, resp, reader := startStream(t, srv, EventStatus)
	defer cancel()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("expected text/event-stream, got %q", got)
	}

	deadline := time.Now().Add(2 * time.Second)
	for e.subscriberCount(EventStatus) != 1 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}

	e.SendJobEvent(true, "echo-route", nil)
	readUntil(t, reader, "echo-route", 2*time.Second)
}

func TestSSEBytesMatchWireFormat(t *testing.T) {
	e := newTestEvent(nil)
	srv := newEventServer(t, e)

	cancel, resp, reader := startStream(t, srv, CommandEvent)
	defer cancel()
	defer resp.Body.Close()

	e.SendCommandEvent(1, "hello")

	frame := readFrame(t, reader, 2*time.Second)

	lines := strings.Split(strings.TrimSuffix(frame, "\n"), "\n")
	if len(lines) != 2 || lines[1] != "" {
		t.Fatalf("expected exactly one data line plus terminator, got %q", frame)
	}
	if !strings.HasPrefix(lines[0], "data: ") {
		t.Fatalf("expected data prefix, got %q", lines[0])
	}

	var info CommandInfo
	if err := json.Unmarshal([]byte(strings.TrimPrefix(lines[0], "data: ")), &info); err != nil {
		t.Fatalf("data line is not valid JSON: %v (%q)", err, lines[0])
	}
	if info.Command != "hello" || info.Severity != 1 {
		t.Fatalf("unexpected payload: %+v", info)
	}
}
