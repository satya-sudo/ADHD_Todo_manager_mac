package ui

import (
	"fmt"

	"focusbar/internal/task"
	"focusbar/internal/tray"

	"github.com/getlantern/systray"
)

type Renderer struct {
	manager        *task.Manager
	onActivity     func()
	currentView    view
	selectedTaskID string
	menuItems      []*systray.MenuItem
}

func New(manager *task.Manager, onActivity func()) *Renderer {
	return &Renderer{
		manager:     manager,
		onActivity:  onActivity,
		currentView: mainView,
	}
}

func (r *Renderer) clearMenu() {
	for _, item := range r.menuItems {
		item.Hide()
	}

	r.menuItems = []*systray.MenuItem{}
}

func (r *Renderer) renderMainMenu() {
	active := r.manager.GetActiveTask()
	if active != nil {
		systray.SetTitle(tray.TaskTitle(active.Title))
	} else {
		systray.SetTitle(tray.IdleTitle())
	}

	if active != nil {
		r.addLabel("Current")
		r.addItem("⚡ " + active.Title)
		r.addAction("⏸ Pause", func() {
			r.manager.PauseTask(active.ID)
			r.goToMain()
		})
		r.addAction("✔ Done", func() {
			r.manager.CompleteTask(active.ID)
			r.goToMain()
		})
	}

	systray.AddSeparator()
	r.renderSection("Todo", task.Todo)
	r.renderSection("Paused", task.Paused)

	systray.AddSeparator()
	r.addAction("➕ Add Task", func() {
		title := promptTask()
		if title == "" {
			return
		}

		r.manager.AddTask(title)
		r.Render()
	})

	systray.AddSeparator()
	r.addAction("Quit", systray.Quit)
}

func (r *Renderer) renderSection(title string, state task.State) {
	filtered := task.FilterByState(r.manager.GetTasks(), state)
	if len(filtered) == 0 {
		return
	}

	r.addLabel(fmt.Sprintf("%s (%d)", title, len(filtered)))
	for _, current := range filtered {
		item := current
		r.addAction(r.formatTitle(item), func() {
			r.recordActivity()
			r.selectTask(item.ID)
			systray.SetTitle(tray.TaskTitle(item.Title))
		})
	}
}

func (r *Renderer) renderTaskActionMenu() {
	current := task.FindByID(r.manager.GetTasks(), r.selectedTaskID)
	if current == nil {
		r.goToMain()
		return
	}

	r.addLabel(current.Title)
	systray.AddSeparator()

	switch current.State {
	case task.Todo:
		r.addAction("▶ Start", func() {
			r.manager.StartTask(current.ID)
			r.goToMain()
		})
	case task.Working:
		r.addAction("⏸ Pause", func() {
			r.manager.PauseTask(current.ID)
			r.goToMain()
		})
	case task.Paused:
		r.addAction("▶ Resume", func() {
			r.manager.StartTask(current.ID)
			r.goToMain()
		})
	}

	r.addAction("✔ Done", func() {
		r.manager.CompleteTask(current.ID)
		r.goToMain()
	})
	r.addAction("🗑 Delete", func() {
		r.manager.DeleteTask(current.ID)
		r.goToMain()
	})

	systray.AddSeparator()
	r.addAction("← Back", r.goToMain)
}

func (r *Renderer) addLabel(title string) {
	item := systray.AddMenuItem(title, "")
	item.Disable()
	r.menuItems = append(r.menuItems, item)
}

func (r *Renderer) addItem(title string) {
	item := systray.AddMenuItem(title, "")
	r.menuItems = append(r.menuItems, item)
}

func (r *Renderer) addAction(title string, action func()) {
	item := systray.AddMenuItem(title, "")
	r.menuItems = append(r.menuItems, item)

	go func() {
		<-item.ClickedCh
		r.recordActivity()
		action()
	}()
}

func (r *Renderer) formatTitle(current task.Task) string {
	switch current.State {
	case task.Working:
		return "⚡ " + current.Title
	case task.Paused:
		return "⏸ " + current.Title
	case task.Done:
		return "✔ " + current.Title
	default:
		return "• " + current.Title
	}
}

func (r *Renderer) recordActivity() {
	if r.onActivity != nil {
		r.onActivity()
	}
}
