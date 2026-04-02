package app

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"focusbar/internal/adaptive"
	"focusbar/internal/logx"
	"focusbar/internal/notifier"
	"focusbar/internal/reminder"
	"focusbar/internal/task"
	"focusbar/internal/ui"

	"github.com/getlantern/systray"
)

type App struct {
	manager        *task.Manager
	renderer       *ui.Renderer
	reminderEngine *reminder.Engine
	adaptiveStats  *adaptive.Stats
	lastActivity   time.Time
	lastNotice     time.Time
	pendingNotice  bool
	mu             sync.RWMutex
}

func New() *App {
	logx.Infof("creating app instance")
	instance := &App{}
	dbPath, legacyPath := defaultStoragePaths()
	instance.adaptiveStats = adaptive.NewStats(12)
	manager := task.NewManager(dbPath, legacyPath, systray.SetTitle, instance.touchActivity, instance.RecordTaskEngagement)
	renderer := ui.New(manager, instance.touchActivity)
	instance.manager = manager
	instance.renderer = renderer
	instance.reminderEngine = reminder.New(instance, systray.SetTitle, notifier.New())
	logx.Infof("using database path=%s legacy_json=%s", dbPath, legacyPath)

	return instance
}

func (a *App) OnReady() {
	logx.Infof("app ready start")
	systray.SetTitle("⚡")
	a.touchActivity()
	a.manager.Load()
	a.manager.CleanupOldTasks()
	a.manager.ResumeActiveTask()
	a.renderer.Render()
	a.manager.StartCleanupLoop()
	a.StartReminder()
	logx.Infof("app ready complete")
}

func (a *App) OnExit() {
	logx.Infof("app exiting")
	if a.reminderEngine != nil {
		a.reminderEngine.Stop()
	}
	if err := a.manager.Close(); err != nil {
		logx.Errorf("task manager close failed err=%v", err)
	}
}

func (a *App) StartReminder() {
	if a.reminderEngine != nil {
		logx.Infof("starting reminder engine")
		a.reminderEngine.Start()
	}
}

func (a *App) GetActiveTaskID() string {
	return a.manager.GetActiveTaskID()
}

func (a *App) GetActiveState() task.TaskState {
	return a.manager.GetActiveState()
}

func (a *App) CountTodo() int {
	return a.manager.CountByState(task.Todo)
}

func (a *App) CountPaused() int {
	return a.manager.CountByState(task.Paused)
}

func (a *App) LastActivity() time.Time {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastActivity
}

func (a *App) HasActiveTask() bool {
	return a.manager.HasActiveTask()
}

func (a *App) LastNotification() time.Time {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastNotice
}

func (a *App) RecordNotification() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.pendingNotice {
		a.adaptiveStats.Record(false, time.Now())
		logx.Infof("adaptive notification missed response_rate=%.2f", a.adaptiveStats.ResponseRate(time.Now()))
	}

	a.lastNotice = time.Now()
	a.pendingNotice = true
	logx.Infof("notification recorded at=%s", a.lastNotice.Format(time.RFC3339))
}

func (a *App) touchActivity() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.lastActivity = time.Now()
	logx.Infof("activity touched at=%s", a.lastActivity.Format(time.RFC3339))
}

func (a *App) RecordTaskEngagement() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.pendingNotice {
		return
	}

	success := time.Since(a.lastNotice) <= 3*time.Minute
	a.adaptiveStats.Record(success, time.Now())
	a.pendingNotice = false
	logx.Infof("adaptive engagement recorded success=%t response_rate=%.2f", success, a.adaptiveStats.ResponseRate(time.Now()))
}

func (a *App) ResponseRate(now time.Time) float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.adaptiveStats.ResponseRate(now)
}

func (a *App) IsBestFocusHour(now time.Time) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.adaptiveStats.IsBestHour(now)
}

func (a *App) IsWeakFocusHour(now time.Time) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.adaptiveStats.IsWeakHour(now)
}

func defaultStoragePaths() (string, string) {
	if path := os.Getenv("FOCUSBAR_DB_PATH"); path != "" {
		legacy := os.Getenv("FOCUSBAR_TASKS_PATH")
		if legacy == "" {
			legacy = filepath.Join(filepath.Dir(path), "tasks.json")
		}
		return path, legacy
	}

	if path := os.Getenv("FOCUSBAR_TASKS_PATH"); path != "" {
		if filepath.Ext(path) == ".db" {
			return path, filepath.Join(filepath.Dir(path), "tasks.json")
		}
		return jsonPathToDBPath(path), path
	}

	baseDir := defaultDataDir()
	if baseDir == "." {
		return "focusbar.db", "tasks.json"
	}

	return filepath.Join(baseDir, "focusbar.db"), filepath.Join(baseDir, "tasks.json")
}

func jsonPathToDBPath(path string) string {
	ext := filepath.Ext(path)
	if ext == ".json" {
		return path[:len(path)-len(ext)] + ".db"
	}

	return path + ".db"
}

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		logx.Errorf("user home lookup failed err=%v", err)
		return "."
	}

	return filepath.Join(home, "Library", "Application Support", "Focusbar")
}
