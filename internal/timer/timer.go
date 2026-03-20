package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	stop      chan struct{}
	startedAt time.Time
}

func (t *Timer) Start(title string, update func(string)) {
	t.Stop()

	t.stop = make(chan struct{})
	t.startedAt = time.Now()

	go func(stop chan struct{}) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				elapsed := time.Since(t.startedAt)
				hours := int(elapsed.Hours())
				minutes := int(elapsed.Minutes()) % 60
				seconds := int(elapsed.Seconds()) % 60
				update(fmt.Sprintf("⚡ %02d:%02d:%02d %s", hours, minutes, seconds, title))
			}
		}
	}(t.stop)
}

func (t *Timer) Stop() {
	if t.stop == nil {
		return
	}

	close(t.stop)
	t.stop = nil
}
