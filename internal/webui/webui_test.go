package webui

import (
	"testing"
)

func TestDistFS_ReturnsFalseWithoutBuild(t *testing.T) {
	// dist/.keep exists but no index.html, so DistFS should return false
	_, ok := DistFS()
	if ok {
		t.Fatal("expected DistFS to return false when dist has no index.html")
	}
}
