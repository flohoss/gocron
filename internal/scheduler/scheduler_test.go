package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestNew_StartsAndStops(t *testing.T) {
	s := New(time.UTC)
	defer s.Stop()

	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

// The location has to reach robfig/cron, otherwise a time_zone change would
// silently keep running schedules in the old zone.
func TestNew_UsesConfiguredLocation(t *testing.T) {
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("failed to load test location: %v", err)
	}

	s := New(berlin)
	defer s.Stop()

	if got := s.Location(); got == nil || got.String() != berlin.String() {
		t.Fatalf("expected scheduler location %q, got %v", berlin, got)
	}
}

func TestAdd_RunsCommandOnSchedule(t *testing.T) {
	s := New(time.UTC)
	defer s.Stop()

	fired := make(chan struct{}, 1)
	s.Add("@every 100ms", func() {
		select {
		case fired <- struct{}{}:
		default:
		}
	})

	select {
	case <-fired:
	case <-time.After(2 * time.Second):
		t.Fatal("scheduled command did not run")
	}
}

func TestAdd_RejectsInvalidCron(t *testing.T) {
	s := New(time.UTC)
	defer s.Stop()

	err := s.Add("not a cron", func() {})
	if err == nil {
		t.Fatal("expected an error for an invalid cron expression")
	}
}

func TestAdd_AcceptsFiveFieldCron(t *testing.T) {
	s := New(time.UTC)
	defer s.Stop()

	if err := s.Add("*/1 * * * *", func() {}); err != nil {
		t.Fatalf("expected five-field cron to be accepted, got %v", err)
	}
}

func TestStop_WaitsForRunningCommand(t *testing.T) {
	s := New(time.UTC)

	started := make(chan struct{})
	released := make(chan struct{})
	var startedOnce sync.Once
	s.Add("@every 100ms", func() {
		startedOnce.Do(func() { close(started) })
		<-released
	})

	<-started

	ctx := s.Stop()
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		t.Fatal("Stop returned before the running command finished")
	case <-timer.C:
	}

	close(released)

	done := time.After(2 * time.Second)
	select {
	case <-ctx.Done():
	case <-done:
		t.Fatal("Stop context never completed after the command finished")
	}
}
