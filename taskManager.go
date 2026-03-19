package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"time"
)

type TaskManager struct {
	tasks          []Task
	activeTaskID   string
	timerStop      chan bool
	focusStartTime time.Time
}

func (m *TaskManager) AddTask(title string) {

	if title == "" {
		return
	}

	task := NewTask(title)

	m.tasks = append(m.tasks, task)

	m.save()
}

func (m *TaskManager) StartTask(id string) {

	// stop existing timer
	m.stopTimer()

	for i := range m.tasks {

		if m.tasks[i].ID == id {

			m.tasks[i].State = Working
			m.activeTaskID = id

		} else if m.tasks[i].State == Working {

			m.tasks[i].State = Paused
		}
	}

	// start new timer
	task := FindTaskByID(m.tasks, id)
	if task != nil {
		m.startTimer(task.Title)
	}

	m.save()
}

func (m *TaskManager) PauseTask(id string) {

	for i := range m.tasks {

		if m.tasks[i].ID == id {

			m.tasks[i].State = Paused

			// stop timer ONLY if this was active
			if m.activeTaskID == id {
				m.stopTimer()
				m.activeTaskID = ""
			}
		}
	}

	m.save()
}

func (m *TaskManager) CompleteTask(id string) {

	for i := range m.tasks {

		if m.tasks[i].ID == id {

			m.tasks[i].State = Done

			if m.activeTaskID == id {
				m.stopTimer()
				m.activeTaskID = ""
			}
		}
	}

	m.save()
}

func (m *TaskManager) DeleteTask(id string) {

	var updated []Task

	for _, t := range m.tasks {

		if t.ID == id {

			if m.activeTaskID == id {
				m.stopTimer()
				m.activeTaskID = ""
			}

			continue
		}

		updated = append(updated, t)
	}

	m.tasks = updated

	m.save()
}

func (m *TaskManager) GetActiveTask() *Task {

	if m.activeTaskID == "" {
		return nil
	}

	return FindTaskByID(m.tasks, m.activeTaskID)
}

func (m *TaskManager) GetTasks() []Task {
	return m.tasks
}

func (m *TaskManager) startTimer(taskTitle string) {

	m.stopTimer()

	m.timerStop = make(chan bool)
	m.focusStartTime = time.Now()

	go func(title string, stop chan bool) {

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {

			case <-stop:
				return

			case <-ticker.C:

				elapsed := time.Since(m.focusStartTime)

				hours := int(elapsed.Hours())
				minutes := int(elapsed.Minutes()) % 60
				seconds := int(elapsed.Seconds()) % 60

				systray.SetTitle(fmt.Sprintf("⚡ %02d:%02d:%02d %s", hours, minutes, seconds, title))
			}
		}

	}(taskTitle, m.timerStop)
}

func (m *TaskManager) stopTimer() {

	if m.timerStop != nil {
		close(m.timerStop)
		m.timerStop = nil
	}
}
