package main

import (
	"github.com/getlantern/systray"
)

var addTaskItem *systray.MenuItem
var quitItem *systray.MenuItem

func onReady() {

	systray.SetTitle("⚡")

	// menu items
	addTaskItem = systray.AddMenuItem("➕ Add Task", "Add new task")
	systray.AddSeparator()

	quitItem = systray.AddMenuItem("Quit", "Quit")

	// load existing tasks
	loadTasks()
	active := getActiveTask()
	if active != "" {
		startFocusTimer(active)
	}
	// remove expired tasks
	cleanupOldTasks()

	// render menu
	renderMenu(addTaskItem)

	// background services
	startReminder()
	startCleanupLoop()

	go menuEventLoop()
}

func menuEventLoop() {

	for {

		select {

		case <-addTaskItem.ClickedCh:

			title := promptTask()

			if title != "" {

				addTask(title)
				renderMenu(addTaskItem)

			}

		case <-quitItem.ClickedCh:

			systray.Quit()
			return
		}
	}
}

func main() {

	systray.Run(onReady, nil)

}
