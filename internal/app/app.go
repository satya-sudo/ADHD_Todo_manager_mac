package app

import (
	"sync"
	"time"

	"focusbar/internal/logx"
	"focusbar/internal/notifier"
	"focusbar/internal/reminder"
	"focusbar/internal/task"
	"focusbar/internal/ui"

	"github.com/getlantern/systray"
)

const storagePath = "tasks.json"

type App struct {
	manager        *task.Manager
	renderer       *ui.Renderer
	reminderEngine *reminder.Engine
	lastActivity   time.Time
	lastNotice     time.Time
	mu             sync.RWMutex
}

func New() *App {
	logx.Infof("creating app instance")
	instance := &App{}
	manager := task.NewManager(storagePath, systray.SetTitle, instance.touchActivity)
	renderer := ui.New(manager, instance.touchActivity)
	instance.manager = manager
	instance.renderer = renderer
	instance.reminderEngine = reminder.New(instance, systray.SetTitle, notifier.New())

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

	a.lastNotice = time.Now()
	logx.Infof("notification recorded at=%s", a.lastNotice.Format(time.RFC3339))
}

func (a *App) touchActivity() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.lastActivity = time.Now()
	logx.Infof("activity touched at=%s", a.lastActivity.Format(time.RFC3339))
}
