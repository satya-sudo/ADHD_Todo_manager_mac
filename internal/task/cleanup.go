package task

import "time"

func (m *Manager) CleanupOldTasks() {
	now := time.Now().Unix()
	var fresh []Task

	for _, current := range m.tasks {
		if now-current.CreatedAt < 86400 || current.State != Done {
			fresh = append(fresh, current)
		}
	}

	m.tasks = fresh
	m.save()
}

func (m *Manager) StartCleanupLoop() {
	go func() {
		for {
			time.Sleep(time.Hour)
			m.CleanupOldTasks()
		}
	}()
}
