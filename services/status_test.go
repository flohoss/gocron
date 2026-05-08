package services

import "testing"

func TestStatus_Int64(t *testing.T) {
	cases := []struct {
		status Status
		want   int64
	}{
		{Running, 1},
		{Stopped, 2},
		{Finished, 3},
	}
	for _, tc := range cases {
		if got := tc.status.Int64(); got != tc.want {
			t.Errorf("Status(%d).Int64() = %d, want %d", tc.status, got, tc.want)
		}
	}
}

func TestSeverity_Values(t *testing.T) {
	if Debug != 1 {
		t.Errorf("expected Debug = 1, got %d", Debug)
	}
	if Info != 2 {
		t.Errorf("expected Info = 2, got %d", Info)
	}
	if Warning != 3 {
		t.Errorf("expected Warning = 3, got %d", Warning)
	}
	if Error != 4 {
		t.Errorf("expected Error = 4, got %d", Error)
	}
}
