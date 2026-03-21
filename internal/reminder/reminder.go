package reminder

import (
	"math/rand"
	"time"

	"focusbar/internal/logx"
	"focusbar/internal/notifier"
	"focusbar/internal/task"
	"focusbar/internal/tray"
)

type Snapshot struct {
	ActiveTaskID string
	ActiveState  task.TaskState

	TodoCount   int
	PausedCount int

	LastActivity     time.Time
	LastNotification time.Time
	IdleDuration     time.Duration
	Now              time.Time
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
	TrayText string
	Message  string
	Priority int
	Notify   bool
}

type Provider interface {
	GetActiveTaskID() string
	GetActiveState() task.TaskState
	CountTodo() int
	CountPaused() int
	LastActivity() time.Time
	LastNotification() time.Time
	RecordNotification()
	HasActiveTask() bool
}

type Engine struct {
	provider Provider
	setTitle func(string)
	notifier notifier.Sender
	stop     chan struct{}
}

func New(provider Provider, setTitle func(string), sender notifier.Sender) *Engine {
	return &Engine{
		provider: provider,
		setTitle: setTitle,
		notifier: sender,
		stop:     make(chan struct{}),
	}
}

func (e *Engine) Start() {
	logx.Infof("reminder engine start ticker=30s")
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
	logx.Infof("reminder engine stop")
	select {
	case <-e.stop:
		return
	default:
		close(e.stop)
	}
}

func (e *Engine) tick() {
	if e.provider.HasActiveTask() {
		logx.Infof("reminder tick skipped active_task=true")
		return
	}

	snapshot := e.buildSnapshot()
	logx.Infof(
		"reminder tick idle=%s todo=%d paused=%d last_activity=%s last_notification=%s",
		snapshot.IdleDuration,
		snapshot.TodoCount,
		snapshot.PausedCount,
		snapshot.LastActivity.Format(time.RFC3339),
		snapshot.LastNotification.Format(time.RFC3339),
	)

	decision := Evaluate(snapshot)
	logx.Infof(
		"reminder decision type=%d tray=%q notify=%t message=%q",
		decision.Type,
		decision.TrayText,
		decision.Notify,
		decision.Message,
	)
	e.apply(decision)
}

func (e *Engine) buildSnapshot() Snapshot {
	now := time.Now()
	lastActivity := e.provider.LastActivity()
	idleDuration := time.Duration(0)
	if !lastActivity.IsZero() {
		idleDuration = now.Sub(lastActivity)
	}

	return Snapshot{
		ActiveTaskID:     e.provider.GetActiveTaskID(),
		ActiveState:      e.provider.GetActiveState(),
		TodoCount:        e.provider.CountTodo(),
		PausedCount:      e.provider.CountPaused(),
		LastActivity:     lastActivity,
		LastNotification: e.provider.LastNotification(),
		IdleDuration:     idleDuration,
		Now:              now,
	}
}

func Evaluate(s Snapshot) Decision {
	if s.ActiveState == task.Working {
		return none()
	}

	if !hasPendingTasks(s) {
		return none()
	}

	var decisions []Decision

	if isIdle(s) {
		decisions = append(decisions, idleDecision(s))
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

	return s.IdleDuration > 2*time.Minute
}

func hasPendingTasks(s Snapshot) bool {
	return s.TodoCount > 0 || s.PausedCount > 0
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

func idleDecision(s Snapshot) Decision {
	if s.IdleDuration > 2*time.Minute && canNotify(s) {
		return Decision{
			Type:     Idle,
			TrayText: tray.NudgeTitle(),
			Message:  "⚡ You've been idle - start one small task",
			Priority: 3,
			Notify:   true,
		}
	}

	if s.IdleDuration > 1*time.Minute {
		return Decision{
			Type:     Idle,
			TrayText: tray.NudgeTitle(),
			Message:  "👀 Still waiting - pick something small",
			Priority: 3,
		}
	}

	messages := []string{
		"👀 Ready when you are",
		"⚡ Pick one thing",
		"🌱 Start small",
	}

	return Decision{
		Type:     Idle,
		TrayText: tray.NudgeTitle(),
		Message:  messages[rand.Intn(len(messages))],
		Priority: 3,
	}
}

func canNotify(s Snapshot) bool {
	if s.LastNotification.IsZero() {
		return true
	}

	return s.Now.Sub(s.LastNotification) > 2*time.Minute
}

func pausedDecision() Decision {
	return Decision{
		Type:     NudgePaused,
		TrayText: tray.ResumeTitle(),
		Message:  "⏸ Resume one task",
		Priority: 2,
	}
}

func overloadDecision() Decision {
	return Decision{
		Type:     Overload,
		TrayText: tray.OverloadTitle(),
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
		logx.Infof("reminder apply skipped type=%d", decision.Type)
		return
	}

	if decision.TrayText != "" {
		e.setTitle(decision.TrayText)
		logx.Infof("reminder tray title updated=%q", decision.TrayText)
	}

	if !decision.Notify {
		logx.Infof("reminder notify skipped notify=false")
		return
	}

	if e.notifier == nil {
		logx.Errorf("reminder notify skipped notifier=nil")
		return
	}

	if err := e.notifier.Notify("Focusbar", decision.Message); err == nil {
		e.provider.RecordNotification()
		return
	}

	logx.Errorf("reminder notify failed message=%q", decision.Message)
}
