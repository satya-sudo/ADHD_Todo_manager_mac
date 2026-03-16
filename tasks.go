package main

import "time"

type TaskState string

const (
	Todo    TaskState = "todo"
	Working TaskState = "working"
	Paused  TaskState = "paused"
	Done    TaskState = "done"
)

type Task struct {
	Title     string    `json:"title"`
	State     TaskState `json:"state"`
	CreatedAt int64     `json:"created_at"`
}

func addTask(title string) {

	if title == "" {
		return
	}

	tasks = append(tasks, Task{
		Title:     title,
		State:     Todo,
		CreatedAt: time.Now().Unix(),
	})

	saveTasks()
}

func startWorkingTask(title string) {

	for i := range tasks {

		if tasks[i].Title == title {

			tasks[i].State = Working
			startFocusTimer(tasks[i].Title)

		} else if tasks[i].State == Working {

			tasks[i].State = Paused

		}

	}

	saveTasks()
}

func pauseTask(title string) {

	for i := range tasks {

		if tasks[i].Title == title {

			tasks[i].State = Paused

		}

	}

	saveTasks()
}

func markTaskDone(title string) {

	for i := range tasks {

		if tasks[i].Title == title {

			tasks[i].State = Done

		}

	}

	saveTasks()
}

func getActiveTask() string {

	for _, t := range tasks {

		if t.State == Working {

			return t.Title

		}

	}

	return ""

}
func setTaskState(title string, state TaskState) {

	for i := range tasks {

		if tasks[i].Title == title {

			tasks[i].State = state

			if state == Working {
				startFocusTimer(title)
			}

		}

	}

	saveTasks()
}
