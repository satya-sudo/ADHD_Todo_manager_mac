package app

import (
	"focusbar/internal/task"
	"focusbar/internal/ui"

	"github.com/getlantern/systray"
)

const storagePath = "tasks.json"

type App struct {
	manager  *task.Manager
	renderer *ui.Renderer
}

func New() *App {
	manager := task.NewManager(storagePath, systray.SetTitle)

	return &App{
		manager:  manager,
		renderer: ui.New(manager),
	}
}

func (a *App) OnReady() {
	systray.SetTitle("⚡")
	a.manager.Load()
	a.manager.CleanupOldTasks()
	a.manager.ResumeActiveTask()
	a.renderer.Render()
	a.manager.StartCleanupLoop()
}

func (a *App) OnExit() {}
