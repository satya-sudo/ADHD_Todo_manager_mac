package adaptive

type Stats struct {
	recent     []bool
	maxEntries int
}

func NewStats(maxEntries int) *Stats {
	if maxEntries <= 0 {
		maxEntries = 12
	}

	return &Stats{maxEntries: maxEntries}
}

func (s *Stats) Record(result bool) {
	s.recent = append(s.recent, result)
	if len(s.recent) > s.maxEntries {
		s.recent = s.recent[len(s.recent)-s.maxEntries:]
	}
}

func (s *Stats) ResponseRate() float64 {
	if len(s.recent) == 0 {
		return 0.5
	}

	success := 0
	for _, result := range s.recent {
		if result {
			success++
		}
	}

	return float64(success) / float64(len(s.recent))
}
