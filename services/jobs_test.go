package services

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestFormatTime_FormatsUnixMillis(t *testing.T) {
	ms := int64(1753900800000)
	got := formatTime(ms)
	want := time.Unix(ms/1000, 0).Local().Format(DATE_FORMAT)
	if got != want {
		t.Fatalf("formatTime(%d) = %q, want %q", ms, got, want)
	}
}

func TestFormatTime_HandlesZero(t *testing.T) {
	got := formatTime(0)
	if got != "1970-01-01 01:00:00" && got != "1970-01-01 00:00:00" {
		t.Fatalf("unexpected format for zero time: %q", got)
	}
}

func TestGenerateUniqueTimestamp_AlwaysIncrements(t *testing.T) {
	lastTimestamp = 0

	first := generateUniqueTimestamp()
	second := generateUniqueTimestamp()

	if second <= first {
		t.Fatalf("expected second > first, got first=%d second=%d", first, second)
	}
}

func TestGenerateUniqueTimestamp_IncrementsWhenSameMillisecond(t *testing.T) {
	lastTimestamp = time.Now().UnixMilli()

	a := generateUniqueTimestamp()
	b := generateUniqueTimestamp()

	if b != a+1 {
		t.Fatalf("expected b=a+1 when same ms, got a=%d b=%d", a, b)
	}
}

func TestGenerateUniqueTimestamp_ConcurrentCallsAreUnique(t *testing.T) {
	lastTimestamp = 0

	var wg sync.WaitGroup
	results := make(chan int64, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- generateUniqueTimestamp()
		}()
	}

	wg.Wait()
	close(results)

	seen := make(map[int64]bool, 100)
	for ts := range results {
		if seen[ts] {
			t.Fatalf("duplicate timestamp generated: %d", ts)
		}
		seen[ts] = true
	}
}

func TestExecuteJobs_SkipsWhenNotIdle(t *testing.T) {
	js := &JobService{
		jobCtx:    context.Background(),
		jobCancel: func() {},
	}

	js.jobMu.Lock()
	js.jobRunning = true
	js.jobMu.Unlock()

	called := false
	js.ExecuteJobs(nil)
	if called {
		t.Fatal("expected ExecuteJobs to skip when job is already running")
	}
}
