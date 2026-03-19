package main

import (
	"github.com/getlantern/systray"
)

var manager = &TaskManager{}

func onReady() {

	// Initial title
	systray.SetTitle("⚡")

	// Load tasks via manager
	manager.load()

	// Cleanup expired tasks
	manager.cleanupOldTasks()

	// Resume active task timer if exists
	active := manager.GetActiveTask()
	if active != nil {
		manager.startTimer(active.Title)
	}

	// First render
	render()

	// Background services
	//startReminder()
	manager.startCleanupLoop()
}

func onExit() {
	// optional cleanup later
}

func main() {
	systray.Run(onReady, onExit)
}
