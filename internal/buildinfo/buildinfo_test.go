package buildinfo

import (
	"strings"
	"testing"
)

func TestSummary_WithoutRepoURL(t *testing.T) {
	Version = "v1.2.3"
	BuildTime = "2026-01-01T00:00:00Z"
	RepoURL = ""

	s := Summary()
	if !strings.Contains(s, "v1.2.3") {
		t.Fatalf("expected version in summary, got %q", s)
	}
	if !strings.Contains(s, "2026-01-01T00:00:00Z") {
		t.Fatalf("expected build time in summary, got %q", s)
	}
	if strings.Contains(s, "source") {
		t.Fatalf("expected no source URL in summary, got %q", s)
	}
}

func TestSummary_WithRepoURL(t *testing.T) {
	Version = "v1.2.3"
	BuildTime = "2026-01-01T00:00:00Z"
	RepoURL = "https://github.com/example/repo"

	s := Summary()
	if !strings.Contains(s, "source") {
		t.Fatalf("expected source URL in summary, got %q", s)
	}
	if !strings.Contains(s, "https://github.com/example/repo") {
		t.Fatalf("expected repo URL in summary, got %q", s)
	}
}
