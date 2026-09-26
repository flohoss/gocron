package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	EventStatus  = "status"
	CommandEvent = "command"
)

const (
	subscriberBuffer = 64
	keepAlivePeriod  = 15 * time.Second
)

type Event struct {
	onSubscribe func(streamID string)
	keepAlive   time.Duration

	mu      sync.Mutex
	streams map[string]map[chan []byte]struct{}
}

type EventInfo struct {
	Idle bool `json:"idle"`
	Run  any  `json:"run"`
	Jobs any  `json:"jobs"`
}

type CommandInfo struct {
	Command  string `json:"command"`
	Severity int    `json:"severity"`
}

func New(onSubscribe func(streamID string)) *Event {
	return &Event{
		onSubscribe: onSubscribe,
		keepAlive:   keepAlivePeriod,
		streams: map[string]map[chan []byte]struct{}{
			EventStatus:  {},
			CommandEvent: {},
		},
	}
}

func (e *Event) SendJobEvent(idle bool, run any, jobs any) {
	data, err := json.Marshal(&EventInfo{
		Idle: idle,
		Run:  run,
		Jobs: jobs,
	})
	if err != nil {
		return
	}
	e.publish(EventStatus, data)
}

func (e *Event) SendCommandEvent(severity int, command string) {
	data, err := json.Marshal(CommandInfo{
		Command:  command,
		Severity: severity,
	})
	if err != nil {
		return
	}
	e.publish(CommandEvent, data)
}

func (e *Event) GetHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		streamID := c.QueryParam("stream")
		if !e.hasStream(streamID) {
			return echo.NewHTTPError(http.StatusNotFound, "Stream not found")
		}

		subscriber := e.subscribe(streamID)
		defer e.unsubscribe(streamID, subscriber)

		w := c.Response()
		headers := w.Header()
		headers.Set(echo.HeaderContentType, "text/event-stream")
		headers.Set(echo.HeaderCacheControl, "no-cache")
		headers.Set(echo.HeaderConnection, "keep-alive")
		headers.Set("X-Accel-Buffering", "no")
		w.WriteHeader(http.StatusOK)

		controller := http.NewResponseController(w)
		if err := controller.Flush(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
		}

		keepAlive := time.NewTicker(e.keepAlive)
		defer keepAlive.Stop()

		for {
			select {
			case <-c.Request().Context().Done():
				return nil
			case <-keepAlive.C:
				if _, err := fmt.Fprint(w, ": keep-alive\n\n"); err != nil {
					return nil
				}
				if err := controller.Flush(); err != nil {
					return nil
				}
			case data, open := <-subscriber:
				if !open {
					return nil
				}
				if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
					return nil
				}
				if err := controller.Flush(); err != nil {
					return nil
				}
			}
		}
	}
}

func (e *Event) publish(streamID string, data []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for subscriber := range e.streams[streamID] {
		sendLatestEvent(subscriber, data)
	}
}

func sendLatestEvent(subscriber chan []byte, data []byte) {
	select {
	case subscriber <- data:
		return
	default:
	}

	select {
	case <-subscriber:
	default:
	}

	select {
	case subscriber <- data:
	default:
	}
}

func (e *Event) subscribe(streamID string) chan []byte {
	subscriber := make(chan []byte, subscriberBuffer)

	e.mu.Lock()
	e.streams[streamID][subscriber] = struct{}{}
	e.mu.Unlock()

	if e.onSubscribe != nil {
		e.onSubscribe(streamID)
	}

	return subscriber
}

func (e *Event) unsubscribe(streamID string, subscriber chan []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.streams[streamID][subscriber]; !exists {
		return
	}
	delete(e.streams[streamID], subscriber)
	close(subscriber)
}

func (e *Event) hasStream(streamID string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, exists := e.streams[streamID]
	return exists
}
