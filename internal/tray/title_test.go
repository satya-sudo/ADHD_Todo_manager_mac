package tray

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTaskTitleTruncatesLongTitles(t *testing.T) {
	t.Parallel()

	title := TaskTitle("this is a very long task title that should be truncated")

	if !strings.HasPrefix(title, "⚡ ") {
		t.Fatalf("expected tray prefix, got %q", title)
	}
	if utf8.RuneCountInString(title) > maxTitleRunes+2 {
		t.Fatalf("expected truncated title length <= %d, got %d", maxTitleRunes+2, utf8.RuneCountInString(title))
	}
	if !strings.HasSuffix(title, "…") {
		t.Fatalf("expected ellipsis, got %q", title)
	}
}

func TestTimerTitleTruncatesOnlyTaskPart(t *testing.T) {
	t.Parallel()

	title := TimerTitle("⚡ 00:12:34", "this is a very long task title")

	if !strings.HasPrefix(title, "⚡ 00:12:34 ") {
		t.Fatalf("unexpected timer prefix: %q", title)
	}
	if !strings.HasSuffix(title, "…") {
		t.Fatalf("expected ellipsis, got %q", title)
	}
}

func TestShortTaskTitleRemainsUnchanged(t *testing.T) {
	t.Parallel()

	title := TaskTitle("Write docs")

	if title != "⚡ Write docs" {
		t.Fatalf("unexpected title: %q", title)
	}
}
