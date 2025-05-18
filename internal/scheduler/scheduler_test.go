package scheduler

import (
	"testing"
)

func TestNewScheduler(t *testing.T) {
	s := New(7)
	if s == nil {
		t.Fatal("Expected scheduler instance, got nil")
	}
	if s.DeleteRunsAfterDays != 7 {
		t.Errorf("Expected DeleteRunsAfterDays to be 7, got %d", s.DeleteRunsAfterDays)
	}
	if s.scheduler == nil {
		t.Error("Expected cron scheduler to be initialized")
	}
}

func TestGetParser(t *testing.T) {
	s := New(1)
	parser := s.GetParser()
	if parser == nil {
		t.Fatal("Expected non-nil parser from GetParser")
	}

	_, err := (*parser).Parse("* * * * *")
	if err != nil {
		t.Errorf("Expected valid cron string to parse, got error: %v", err)
	}
}
