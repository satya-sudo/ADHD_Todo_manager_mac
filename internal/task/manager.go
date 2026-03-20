package task

import "focusbar/internal/timer"

type Manager struct {
	tasks        []Task
	activeTaskID string
	storagePath  string
	timer        timer.Timer
	setTitle     func(string)
}

func NewManager(storagePath string, setTitle func(string)) *Manager {
	return &Manager{
		storagePath: storagePath,
		setTitle:    setTitle,
	}
}

func (m *Manager) AddTask(title string) {
	if title == "" {
		return
	}

	m.tasks = append(m.tasks, New(title))
	m.save()
}

func (m *Manager) StartTask(id string) {
	m.timer.Stop()

	for i := range m.tasks {
		switch {
		case m.tasks[i].ID == id:
			m.tasks[i].State = Working
			m.activeTaskID = id
		case m.tasks[i].State == Working:
			m.tasks[i].State = Paused
		}
	}

	if current := FindByID(m.tasks, id); current != nil {
		m.timer.Start(current.Title, m.setTitle)
	}

	m.save()
}

func (m *Manager) PauseTask(id string) {
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

	m.save()
}

func (m *Manager) CompleteTask(id string) {
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

	m.save()
}

func (m *Manager) DeleteTask(id string) {
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
	m.save()
}

func (m *Manager) GetTasks() []Task {
	return m.tasks
}

func (m *Manager) GetActiveTask() *Task {
	if m.activeTaskID == "" {
		return nil
	}

	return FindByID(m.tasks, m.activeTaskID)
}

func (m *Manager) ResumeActiveTask() {
	current := m.GetActiveTask()
	if current == nil {
		return
	}

	m.timer.Start(current.Title, m.setTitle)
}
