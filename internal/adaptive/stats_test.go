package adaptive

import (
	"reflect"
	"testing"
	"time"
)

func TestBestWindowsGroupsAdjacentHours(t *testing.T) {
	t.Parallel()

	stats := NewStats(10)

	recordHour(stats, 9, true, true, true)
	recordHour(stats, 10, true, true, true)
	recordHour(stats, 11, true, true, true)
	recordHour(stats, 15, true, true, true)

	got := stats.BestWindows()
	want := []FocusWindow{
		{StartHour: 9, EndHour: 11},
		{StartHour: 15, EndHour: 15},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected windows: got %#v want %#v", got, want)
	}
}

func TestBestWindowsIgnoresWeakOrInsufficientHours(t *testing.T) {
	t.Parallel()

	stats := NewStats(10)

	recordHour(stats, 9, true, true)
	recordHour(stats, 10, true, true, true)
	recordHour(stats, 11, false, false, false, true)

	got := stats.BestWindows()
	if len(got) != 1 {
		t.Fatalf("expected 1 window, got %d", len(got))
	}
	if got[0] != (FocusWindow{StartHour: 10, EndHour: 10}) {
		t.Fatalf("unexpected window: %#v", got[0])
	}
}

func TestFocusWindowLabel(t *testing.T) {
	t.Parallel()

	window := FocusWindow{StartHour: 10, EndHour: 12}
	if got := window.Label(); got != "10 AM - 1 PM" {
		t.Fatalf("unexpected label: %q", got)
	}

	single := FocusWindow{StartHour: 23, EndHour: 23}
	if got := single.Label(); got != "11 PM - 12 AM" {
		t.Fatalf("unexpected overnight label: %q", got)
	}
}

func recordHour(stats *Stats, hour int, results ...bool) {
	base := time.Date(2026, 3, 26, hour, 0, 0, 0, time.UTC)
	for _, result := range results {
		stats.Record(result, base)
	}
}
