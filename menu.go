package main

import "github.com/getlantern/systray"

var taskItems []*systray.MenuItem

var menuItems []*systray.MenuItem

func renderMenu(addTaskItem *systray.MenuItem) {

	for _, item := range menuItems {
		item.Hide()
	}

	menuItems = []*systray.MenuItem{}

	renderSection("Working", Working)
	renderSection("Todo", Todo)
	renderSection("Paused", Paused)
	renderSection("Done", Done)
}

func renderSection(title string, state TaskState) {

	found := false

	for _, t := range tasks {
		if t.State == state {
			found = true
			break
		}
	}

	if !found {
		return
	}

	header := systray.AddMenuItem(title, "")
	header.Disable()

	menuItems = append(menuItems, header)

	for i := range tasks {

		task := &tasks[i]

		if task.State != state {
			continue
		}

		item := systray.AddMenuItem(formatTitle(*task), "")

		menuItems = append(menuItems, item)

		go handleTaskClick(task, item)
	}

	systray.AddSeparator()
}
func handleTaskClick(task *Task, item *systray.MenuItem) {

	for range item.ClickedCh {

		switch task.State {

		case Todo:
			startWorkingTask(task.Title)

		case Working:
			pauseTask(task.Title)

		case Paused:
			startWorkingTask(task.Title)

		}

		renderMenu(nil)
	}
}
func renderTasks() {

	for _, item := range taskItems {
		item.Hide()
	}

	taskItems = []*systray.MenuItem{}

	if len(tasks) == 0 {
		return
	}

	systray.AddSeparator()

	header := systray.AddMenuItem("Today", "")
	header.Disable()

	for i := range tasks {

		task := &tasks[i]

		item := systray.AddMenuItem("• "+task.Title, "task")

		taskItems = append(taskItems, item)

		go func(t *Task, menuItem *systray.MenuItem) {

			for range menuItem.ClickedCh {

				startWorkingTask(t.Title)

			}

		}(task, item)

	}
}

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
