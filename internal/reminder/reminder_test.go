package reminder

import (
	"testing"
	"time"

	"focusbar/internal/adaptive"
	"focusbar/internal/task"
	"focusbar/internal/tray"
)

func TestEvaluateReturnsNoneWithoutPendingTasks(t *testing.T) {
	t.Parallel()

	decision := Evaluate(Snapshot{
		Now:          time.Now(),
		LastActivity: time.Now().Add(-10 * time.Minute),
		IdleDuration: 10 * time.Minute,
		ResponseRate: 0.5,
		Bucket:       adaptive.Afternoon,
		BestHour:     false,
		WeakHour:     false,
	})

	if decision.Type != None {
		t.Fatalf("expected none, got %v", decision.Type)
	}
}

func TestEvaluateReturnsNoneWhileWorking(t *testing.T) {
	t.Parallel()

	decision := Evaluate(Snapshot{
		ActiveTaskID: "task-1",
		ActiveState:  task.Working,
		TodoCount:    2,
		LastActivity: time.Now().Add(-10 * time.Minute),
		IdleDuration: 10 * time.Minute,
		Now:          time.Now(),
		PausedCount:  0,
		ResponseRate: 0.5,
		Bucket:       adaptive.Afternoon,
		BestHour:     false,
		WeakHour:     false,
	})

	if decision.Type != None {
		t.Fatalf("expected none, got %v", decision.Type)
	}
}

func TestEvaluateIdleEscalatesToNotification(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:        1,
		LastActivity:     now.Add(-13 * time.Minute),
		LastNotification: time.Time{},
		IdleDuration:     13 * time.Minute,
		Now:              now,
		ResponseRate:     0.8,
		Bucket:           adaptive.Afternoon,
		BestHour:         false,
		WeakHour:         false,
	})

	if decision.Type != Idle {
		t.Fatalf("expected idle, got %v", decision.Type)
	}
	if !decision.Notify {
		t.Fatalf("expected notify=true")
	}
	if decision.TrayText != tray.NudgeTitle() {
		t.Fatalf("expected tray text %q, got %q", tray.NudgeTitle(), decision.TrayText)
	}
}

func TestEvaluateIdleRespectsNotificationCooldown(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:        1,
		LastActivity:     now.Add(-13 * time.Minute),
		LastNotification: now.Add(-5 * time.Minute),
		IdleDuration:     13 * time.Minute,
		Now:              now,
		ResponseRate:     0.8,
		Bucket:           adaptive.Afternoon,
		BestHour:         false,
		WeakHour:         false,
	})

	if decision.Type != Idle {
		t.Fatalf("expected idle, got %v", decision.Type)
	}
	if decision.Notify {
		t.Fatalf("expected notify=false during cooldown")
	}
}

func TestEvaluateChoosesPausedOverOverload(t *testing.T) {
	t.Parallel()

	decision := Evaluate(Snapshot{
		TodoCount:    8,
		PausedCount:  3,
		LastActivity: time.Now(),
		IdleDuration: 0,
		Now:          time.Now(),
		ResponseRate: 0.5,
		Bucket:       adaptive.Afternoon,
		BestHour:     false,
		WeakHour:     false,
	})

	if decision.Type != NudgePaused {
		t.Fatalf("expected paused nudge, got %v", decision.Type)
	}
	if decision.TrayText != tray.ResumeTitle() {
		t.Fatalf("expected tray text %q, got %q", tray.ResumeTitle(), decision.TrayText)
	}
}

func TestEvaluateUsesSofterToneForLowResponseRate(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:    1,
		LastActivity: now.Add(-4 * time.Minute),
		IdleDuration: 4 * time.Minute,
		Now:          now,
		ResponseRate: 0.2,
		Bucket:       adaptive.Morning,
		BestHour:     false,
		WeakHour:     false,
	})

	if decision.Type != Idle {
		t.Fatalf("expected idle, got %v", decision.Type)
	}
	if decision.Notify {
		t.Fatalf("expected notify=false")
	}
	valid := map[string]bool{
		"🌱 Start small - no rush":   true,
		"☀️ Ease in with one thing": true,
	}
	if !valid[decision.Message] {
		t.Fatalf("unexpected soft message %q", decision.Message)
	}
}

func TestEvaluateLowResponseRateRequiresLongerIdleForNotification(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:        1,
		LastActivity:     now.Add(-13 * time.Minute),
		LastNotification: time.Time{},
		IdleDuration:     13 * time.Minute,
		Now:              now,
		ResponseRate:     0.2,
		Bucket:           adaptive.Morning,
		BestHour:         false,
		WeakHour:         false,
	})

	if decision.Notify {
		t.Fatalf("expected notify=false for low response rate before longer morning idle")
	}
}

func TestEvaluateReturnsNoneAtNight(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:    2,
		PausedCount:  1,
		LastActivity: now.Add(-30 * time.Minute),
		IdleDuration: 30 * time.Minute,
		Now:          now,
		ResponseRate: 0.8,
		Bucket:       adaptive.Night,
		BestHour:     false,
		WeakHour:     false,
	})

	if decision.Type != None {
		t.Fatalf("expected none at night, got %v", decision.Type)
	}
}

func TestEvaluateBestHourEscalatesSooner(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:        1,
		LastActivity:     now.Add(-8 * time.Minute),
		LastNotification: time.Time{},
		IdleDuration:     8 * time.Minute,
		Now:              now,
		ResponseRate:     0.8,
		Bucket:           adaptive.Afternoon,
		BestHour:         true,
		WeakHour:         false,
	})

	if !decision.Notify {
		t.Fatalf("expected notify=true during best hour")
	}
}

func TestEvaluateWeakHourBacksOff(t *testing.T) {
	t.Parallel()

	now := time.Now()
	decision := Evaluate(Snapshot{
		TodoCount:        1,
		LastActivity:     now.Add(-16 * time.Minute),
		LastNotification: time.Time{},
		IdleDuration:     16 * time.Minute,
		Now:              now,
		ResponseRate:     0.5,
		Bucket:           adaptive.Afternoon,
		BestHour:         false,
		WeakHour:         true,
	})

	if decision.Notify {
		t.Fatalf("expected notify=false during weak hour before longer idle")
	}
}
