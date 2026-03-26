package reminder

import (
	"math/rand"
	"time"

	"focusbar/internal/adaptive"
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
	Bucket       adaptive.TimeBucket
	BestHour     bool
	WeakHour     bool

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
	ResponseRate(now time.Time) float64
	IsBestFocusHour(now time.Time) bool
	IsWeakFocusHour(now time.Time) bool
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
		"reminder tick idle=%s todo=%d paused=%d response_rate=%.2f best_hour=%t weak_hour=%t last_activity=%s last_notification=%s",
		snapshot.IdleDuration,
		snapshot.TodoCount,
		snapshot.PausedCount,
		snapshot.ResponseRate,
		snapshot.BestHour,
		snapshot.WeakHour,
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
		ResponseRate:     e.provider.ResponseRate(now),
		Bucket:           adaptive.GetTimeBucket(now),
		BestHour:         e.provider.IsBestFocusHour(now),
		WeakHour:         e.provider.IsWeakFocusHour(now),
		LastActivity:     lastActivity,
		LastNotification: e.provider.LastNotification(),
		IdleDuration:     idleDuration,
		Now:              now,
	}
}

func Evaluate(s Snapshot) Decision {
	if s.Bucket == adaptive.Night {
		return none()
	}

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
	message := idleMessage(s.ResponseRate, s.Bucket)
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
	if s.WeakHour {
		return s.IdleDuration > 18*time.Minute
	}

	if s.BestHour {
		return s.IdleDuration > 7*time.Minute
	}

	if s.Bucket == adaptive.Morning && s.ResponseRate < 0.3 {
		return s.IdleDuration > 18*time.Minute
	}

	if s.Bucket == adaptive.Afternoon && s.ResponseRate > 0.7 {
		return s.IdleDuration > 8*time.Minute
	}

	if s.Bucket == adaptive.Evening {
		return s.IdleDuration > 14*time.Minute
	}

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
		return s.Bucket != adaptive.Night
	}

	if s.Bucket == adaptive.Night {
		return false
	}

	sinceNotification := s.Now.Sub(s.LastNotification)

	if s.WeakHour {
		return sinceNotification > 22*time.Minute
	}

	if s.BestHour {
		return sinceNotification > 7*time.Minute
	}

	if s.Bucket == adaptive.Morning && responseRate < 0.3 {
		return sinceNotification > 25*time.Minute
	}

	if s.Bucket == adaptive.Afternoon && responseRate > 0.7 {
		return sinceNotification > 8*time.Minute
	}

	if s.Bucket == adaptive.Evening {
		return sinceNotification > 18*time.Minute
	}

	if responseRate < 0.3 {
		return sinceNotification > 20*time.Minute
	}

	if responseRate > 0.7 {
		return sinceNotification > 10*time.Minute
	}

	return sinceNotification > 15*time.Minute
}

func idleMessage(responseRate float64, bucket adaptive.TimeBucket) string {
	if bucket != adaptive.Night && responseRate > 0.7 {
		switch bucket {
		case adaptive.Morning:
			messages := []string{
				"⚡ Morning focus - pick one thing",
				"☀️ You are rolling - start here",
			}

			return messages[rand.Intn(len(messages))]
		case adaptive.Afternoon:
			messages := []string{
				"⚡ This is your window - start now",
				"🔥 Good window - start now",
			}

			return messages[rand.Intn(len(messages))]
		}
	}

	if bucket != adaptive.Night && responseRate < 0.3 {
		switch bucket {
		case adaptive.Morning:
			messages := []string{
				"🌱 Start small - no rush",
				"☀️ Ease in with one thing",
			}

			return messages[rand.Intn(len(messages))]
		case adaptive.Afternoon:
			messages := []string{
				"👀 One small win is enough",
				"🌱 Try one task when ready",
			}

			return messages[rand.Intn(len(messages))]
		}
	}

	switch bucket {
	case adaptive.Morning:
		messages := []string{
			"🌱 Start small - no rush",
			"☀️ Ease in with one thing",
		}

		return messages[rand.Intn(len(messages))]
	case adaptive.Afternoon:
		messages := []string{
			"⚡ Pick one thing",
			"👀 Ready when you are",
		}

		return messages[rand.Intn(len(messages))]
	case adaptive.Evening:
		messages := []string{
			"👀 Maybe one small thing?",
			"🌙 A small step still counts",
		}

		return messages[rand.Intn(len(messages))]
	default:
		return ""
	}
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
