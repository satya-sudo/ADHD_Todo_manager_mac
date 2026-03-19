package main

type View string

const (
	MainView       View = "main"
	TaskActionView View = "task_actions"
)

var currentView View = MainView
var selectedTaskID string

// -------------------- NAVIGATION --------------------

func goToMain() {
	currentView = MainView
	selectedTaskID = ""
	render()
}

func selectTask(id string) {
	selectedTaskID = id
	currentView = TaskActionView
	render()
}

// -------------------- RENDER ENTRY --------------------

func render() {

	clearMenu()

	switch currentView {

	case MainView:
		renderMainMenu()

	case TaskActionView:
		renderTaskActionMenu()
	}
}
