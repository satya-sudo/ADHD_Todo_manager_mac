package reminder

import (
	"math/rand"
	"time"

	"focusbar/internal/task"
)

type Snapshot struct {
	ActiveTaskID string
	ActiveState  task.TaskState

	TodoCount   int
	PausedCount int

	LastActivity time.Time
	Now          time.Time
}

type Type int

const (
	None Type = iota
	Idle
	NudgePaused
	Overload
)

type Decision struct {
	Type     Type
	Message  string
	Priority int
}

type Provider interface {
	GetActiveTaskID() string
	GetActiveState() task.TaskState
	CountTodo() int
	CountPaused() int
	LastActivity() time.Time
	HasActiveTask() bool
}

type Engine struct {
	provider Provider
	setTitle func(string)
	stop     chan struct{}
}

func New(provider Provider, setTitle func(string)) *Engine {
	return &Engine{
		provider: provider,
		setTitle: setTitle,
		stop:     make(chan struct{}),
	}
}

func (e *Engine) Start() {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				e.tick()
			case <-e.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (e *Engine) Stop() {
	select {
	case <-e.stop:
		return
	default:
		close(e.stop)
	}
}

func (e *Engine) tick() {
	if e.provider.HasActiveTask() {
		return
	}

	decision := Evaluate(e.buildSnapshot())
	e.apply(decision)
}

func (e *Engine) buildSnapshot() Snapshot {
	return Snapshot{
		ActiveTaskID: e.provider.GetActiveTaskID(),
		ActiveState:  e.provider.GetActiveState(),
		TodoCount:    e.provider.CountTodo(),
		PausedCount:  e.provider.CountPaused(),
		LastActivity: e.provider.LastActivity(),
		Now:          time.Now(),
	}
}

func Evaluate(s Snapshot) Decision {
	if s.ActiveState == task.Working {
		return none()
	}

	var decisions []Decision

	if isIdle(s) {
		decisions = append(decisions, idleDecision())
	}

	if tooManyPaused(s) {
		decisions = append(decisions, pausedDecision())
	}

	if overload(s) {
		decisions = append(decisions, overloadDecision())
	}

	if len(decisions) == 0 {
		return none()
	}

	return pickHighest(decisions)
}

func isIdle(s Snapshot) bool {
	if s.ActiveTaskID != "" {
		return false
	}

	if s.LastActivity.IsZero() {
		return false
	}

	return s.Now.Sub(s.LastActivity) > 2*time.Minute
}

func tooManyPaused(s Snapshot) bool {
	return s.PausedCount >= 3
}

func overload(s Snapshot) bool {
	return s.TodoCount > 7
}

func none() Decision {
	return Decision{Type: None}
}

func idleDecision() Decision {
	messages := []string{
		"👀 Ready when you are",
		"⚡ Pick one thing",
		"🌱 Start small",
	}

	return Decision{
		Type:     Idle,
		Message:  messages[rand.Intn(len(messages))],
		Priority: 3,
	}
}

func pausedDecision() Decision {
	return Decision{
		Type:     NudgePaused,
		Message:  "⏸ Resume one task",
		Priority: 2,
	}
}

func overloadDecision() Decision {
	return Decision{
		Type:     Overload,
		Message:  "📋 Too many tasks - pick 1",
		Priority: 1,
	}
}

func pickHighest(decisions []Decision) Decision {
	best := decisions[0]

	for _, current := range decisions[1:] {
		if current.Priority > best.Priority {
			best = current
		}
	}

	return best
}

func (e *Engine) apply(decision Decision) {
	if decision.Type == None || e.setTitle == nil {
		return
	}

	e.setTitle(decision.Message)
}
