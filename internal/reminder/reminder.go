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

	TodoCount    int
	PausedCount  int
	ResponseRate float64

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
	ResponseRate() float64
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
		"reminder tick idle=%s todo=%d paused=%d response_rate=%.2f last_activity=%s last_notification=%s",
		snapshot.IdleDuration,
		snapshot.TodoCount,
		snapshot.PausedCount,
		snapshot.ResponseRate,
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
		ResponseRate:     e.provider.ResponseRate(),
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
	message := idleMessage(s.ResponseRate)
	if shouldEscalateIdle(s) && canNotify(s, s.ResponseRate) {
		return Decision{
			Type:     Idle,
			TrayText: tray.NudgeTitle(),
			Message:  message,
			Priority: 3,
			Notify:   true,
		}
	}

	return Decision{
		Type:     Idle,
		TrayText: tray.NudgeTitle(),
		Message:  message,
		Priority: 3,
	}
}

func shouldEscalateIdle(s Snapshot) bool {
	if s.ResponseRate < 0.3 {
		return s.IdleDuration > 15*time.Minute
	}

	if s.ResponseRate > 0.7 {
		return s.IdleDuration > 10*time.Minute
	}

	return s.IdleDuration > 12*time.Minute
}

func canNotify(s Snapshot, responseRate float64) bool {
	if s.LastNotification.IsZero() {
		return true
	}

	sinceNotification := s.Now.Sub(s.LastNotification)

	if responseRate < 0.3 {
		return sinceNotification > 20*time.Minute
	}

	if responseRate > 0.7 {
		return sinceNotification > 10*time.Minute
	}

	return sinceNotification > 15*time.Minute
}

func idleMessage(responseRate float64) string {
	if responseRate < 0.3 {
		messages := []string{
			"🌱 No rush - start small when ready",
			"👀 It is okay, just begin with one thing",
		}

		return messages[rand.Intn(len(messages))]
	}

	if responseRate > 0.7 {
		messages := []string{
			"⚡ Let us go - pick one thing",
			"🔥 You have got this - start now",
		}

		return messages[rand.Intn(len(messages))]
	}

	messages := []string{
		"👀 Ready when you are",
		"⚡ Pick one thing",
		"🌱 Start small",
	}

	return messages[rand.Intn(len(messages))]
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
