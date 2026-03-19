package main

import (
	"fmt"

	"github.com/getlantern/systray"
)

var menuItems []*systray.MenuItem

// -------------------- CORE --------------------

func clearMenu() {
	for _, item := range menuItems {
		item.Hide()
	}
	menuItems = []*systray.MenuItem{}
}

// -------------------- MAIN MENU --------------------

func renderMainMenu() {

	active := manager.GetActiveTask()

	if active != nil {
		systray.SetTitle("⚡ " + active.Title)
	} else {
		systray.SetTitle("⚡ Idle")
	}

	// -------- Current --------
	if active != nil {

		addLabel("Current")
		addItem("⚡ " + active.Title)

		addAction("⏸ Pause", func() {
			manager.PauseTask(active.ID)
			goToMain()
		})

		addAction("✔ Done", func() {
			manager.CompleteTask(active.ID)
			goToMain()
		})
	}

	systray.AddSeparator()

	renderSection("Todo", Todo)
	renderSection("Paused", Paused)

	systray.AddSeparator()

	addAction("➕ Add Task", func() {
		title := promptTask()
		if title != "" {
			manager.AddTask(title)
			render()
		}
	})

	systray.AddSeparator()

	addAction("Quit", func() {
		systray.Quit()
	})
}

// -------------------- TASK SECTION --------------------

func renderSection(title string, state TaskState) {

	var filtered []Task

	for _, t := range manager.GetTasks() {
		if t.State == state {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) == 0 {
		return
	}

	addLabel(fmt.Sprintf("%s (%d)", title, len(filtered)))

	for _, task := range filtered {

		t := task // capture safely

		addAction(formatTitle(t), func() {
			selectTask(t.ID)
			systray.SetTitle("⚡ " + t.Title)
		})
	}
}

// -------------------- TASK ACTION MENU --------------------

func renderTaskActionMenu() {

	task := FindTaskByID(manager.GetTasks(), selectedTaskID)

	if task == nil {
		goToMain()
		return
	}

	addLabel(task.Title)

	systray.AddSeparator()

	switch task.State {

	case Todo:

		addAction("▶ Start", func() {
			manager.StartTask(task.ID)
			goToMain()
		})

	case Working:

		addAction("⏸ Pause", func() {
			manager.PauseTask(task.ID)
			goToMain()
		})

	case Paused:

		addAction("▶ Resume", func() {
			manager.StartTask(task.ID)
			goToMain()
		})
	}

	addAction("✔ Done", func() {
		manager.CompleteTask(task.ID)
		goToMain()
	})

	addAction("🗑 Delete", func() {
		manager.DeleteTask(task.ID)
		goToMain()
	})

	systray.AddSeparator()

	addAction("← Back", func() {
		goToMain()
	})
}

// -------------------- HELPERS --------------------

func addLabel(title string) {
	item := systray.AddMenuItem(title, "")
	item.Disable()
	menuItems = append(menuItems, item)
}

func addItem(title string) {
	item := systray.AddMenuItem(title, "")
	menuItems = append(menuItems, item)
}

func addAction(title string, action func()) {

	item := systray.AddMenuItem(title, "")
	menuItems = append(menuItems, item)

	go func() {
		<-item.ClickedCh
		action()
	}()
}

// -------------------- FORMAT --------------------

func formatTitle(t Task) string {

	switch t.State {

	case Working:
		return "⚡ " + t.Title

	case Paused:
		return "⏸ " + t.Title

	case Done:
		return "✔ " + t.Title

	default:
		return "• " + t.Title
	}
}
