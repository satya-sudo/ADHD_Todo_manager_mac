package adaptive

import (
	"fmt"
	"sync"
	"time"
)

type TimeBucket int

const (
	Morning TimeBucket = iota
	Afternoon
	Evening
	Night
)

type BucketStats struct {
	recent []bool
}

type HourStats struct {
	success int
	total   int
}

type FocusWindow struct {
	StartHour int
	EndHour   int
}

type Stats struct {
	mu      sync.Mutex
	buckets map[TimeBucket]*BucketStats
	hours   [24]HourStats
	max     int
}

func NewStats(max int) *Stats {
	if max <= 0 {
		max = 10
	}

	return &Stats{
		max: max,
		buckets: map[TimeBucket]*BucketStats{
			Morning:   {},
			Afternoon: {},
			Evening:   {},
			Night:     {},
		},
	}
}

func GetTimeBucket(now time.Time) TimeBucket {
	hour := now.Hour()

	switch {
	case hour >= 6 && hour < 12:
		return Morning
	case hour >= 12 && hour < 18:
		return Afternoon
	case hour >= 18 && hour < 23:
		return Evening
	default:
		return Night
	}
}

func (s *Stats) Record(success bool, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket := GetTimeBucket(now)
	stats := s.buckets[bucket]
	stats.recent = append(stats.recent, success)

	if len(stats.recent) > s.max {
		stats.recent = stats.recent[len(stats.recent)-s.max:]
	}

	hour := now.Hour()
	s.hours[hour].total++
	if success {
		s.hours[hour].success++
	}
}

func (s *Stats) ResponseRate(now time.Time) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := s.buckets[GetTimeBucket(now)].recent
	if len(data) == 0 {
		return 0.5
	}

	success := 0
	for _, result := range data {
		if result {
			success++
		}
	}

	return float64(success) / float64(len(data))
}

func (s *Stats) Snapshot(now time.Time) []bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := s.buckets[GetTimeBucket(now)].recent
	out := make([]bool, len(data))
	copy(out, data)
	return out
}

func (s *Stats) HourScore(hour int) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.hourScoreLocked(hour)
}

func (s *Stats) BestHours() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.bestHoursLocked()
}

func (s *Stats) IsBestHour(now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	hour := now.Hour()
	return s.hours[hour].total >= 3 && s.hourScoreLocked(hour) > 0.7
}

func (s *Stats) IsWeakHour(now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	hour := now.Hour()
	return s.hours[hour].total >= 3 && s.hourScoreLocked(hour) < 0.3
}

func (s *Stats) hourScoreLocked(hour int) float64 {
	if hour < 0 || hour > 23 {
		return 0.5
	}

	stats := s.hours[hour]
	if stats.total == 0 {
		return 0.5
	}

	return float64(stats.success) / float64(stats.total)
}

func (s *Stats) BestWindows() []FocusWindow {
	s.mu.Lock()
	defer s.mu.Unlock()

	hours := s.bestHoursLocked()
	if len(hours) == 0 {
		return nil
	}

	windows := make([]FocusWindow, 0, len(hours))
	current := FocusWindow{
		StartHour: hours[0],
		EndHour:   hours[0],
	}

	for _, hour := range hours[1:] {
		if hour == current.EndHour+1 {
			current.EndHour = hour
			continue
		}

		windows = append(windows, current)
		current = FocusWindow{
			StartHour: hour,
			EndHour:   hour,
		}
	}

	windows = append(windows, current)
	return windows
}

func (w FocusWindow) Label() string {
	start := formatHour(w.StartHour)
	end := formatHour((w.EndHour + 1) % 24)

	if w.StartHour == w.EndHour {
		return fmt.Sprintf("%s - %s", start, end)
	}

	return fmt.Sprintf("%s - %s", start, end)
}

func formatHour(hour int) string {
	h := hour % 24
	if h < 0 {
		h += 24
	}

	suffix := "AM"
	display := h

	switch {
	case h == 0:
		display = 12
	case h == 12:
		suffix = "PM"
		display = 12
	case h > 12:
		suffix = "PM"
		display = h - 12
	}

	if h > 0 && h < 12 {
		suffix = "AM"
	}

	return fmt.Sprintf("%d %s", display, suffix)
}

func (s *Stats) bestHoursLocked() []int {
	var best []int
	for hour := 0; hour < 24; hour++ {
		if s.hours[hour].total >= 3 && s.hourScoreLocked(hour) > 0.7 {
			best = append(best, hour)
		}
	}

	return best
}
