package task

import (
	"sync"

	"focusbar/internal/timer"
)

type Manager struct {
	tasks        []Task
	activeTaskID string
	storagePath  string
	timer        timer.Timer
	setTitle     func(string)
	touch        func()
	mu           sync.RWMutex
}

func NewManager(storagePath string, setTitle func(string), touch func()) *Manager {
	return &Manager{
		storagePath: storagePath,
		setTitle:    setTitle,
		touch:       touch,
	}
}

func (m *Manager) AddTask(title string) {
	if title == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks = append(m.tasks, New(title))
	m.recordActivity()
	m.saveLocked()
}

func (m *Manager) StartTask(id string) {
	m.timer.Stop()

	m.mu.Lock()
	defer m.mu.Unlock()

	var title string
	for i := range m.tasks {
		switch {
		case m.tasks[i].ID == id:
			m.tasks[i].State = Working
			m.activeTaskID = id
			title = m.tasks[i].Title
		case m.tasks[i].State == Working:
			m.tasks[i].State = Paused
		}
	}

	if title != "" {
		m.timer.Start(title, m.setTitle)
	}

	m.recordActivity()
	m.saveLocked()
}

func (m *Manager) PauseTask(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.tasks {
		if m.tasks[i].ID != id {
			continue
		}

		m.tasks[i].State = Paused
		if m.activeTaskID == id {
			m.timer.Stop()
			m.activeTaskID = ""
		}
	}

	m.recordActivity()
	m.saveLocked()
}

func (m *Manager) CompleteTask(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.tasks {
		if m.tasks[i].ID != id {
			continue
		}

		m.tasks[i].State = Done
		if m.activeTaskID == id {
			m.timer.Stop()
			m.activeTaskID = ""
		}
	}

	m.recordActivity()
	m.saveLocked()
}

func (m *Manager) DeleteTask(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var updated []Task

	for _, current := range m.tasks {
		if current.ID == id {
			if m.activeTaskID == id {
				m.timer.Stop()
				m.activeTaskID = ""
			}
			continue
		}

		updated = append(updated, current)
	}

	m.tasks = updated
	m.recordActivity()
	m.saveLocked()
}

func (m *Manager) GetTasks() []Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]Task, len(m.tasks))
	copy(tasks, m.tasks)
	return tasks
}

func (m *Manager) GetActiveTask() *Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeTaskID == "" {
		return nil
	}

	current := FindByID(m.tasks, m.activeTaskID)
	if current == nil {
		return nil
	}

	copyTask := *current
	return &copyTask
}

func (m *Manager) GetActiveTaskID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.activeTaskID
}

func (m *Manager) GetActiveState() State {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeTaskID == "" {
		return ""
	}

	current := FindByID(m.tasks, m.activeTaskID)
	if current == nil {
		return ""
	}

	return current.State
}

func (m *Manager) CountByState(state State) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, current := range m.tasks {
		if current.State == state {
			count++
		}
	}

	return count
}

func (m *Manager) HasActiveTask() bool {
	return m.GetActiveTaskID() != ""
}

func (m *Manager) ResumeActiveTask() {
	current := m.GetActiveTask()
	if current == nil {
		return
	}

	m.timer.Start(current.Title, m.setTitle)
}

func (m *Manager) recordActivity() {
	if m.touch != nil {
		m.touch()
	}
}
