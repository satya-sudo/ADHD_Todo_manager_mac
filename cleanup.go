package main

import "time"

func (m *TaskManager) cleanupOldTasks() {

	now := time.Now().Unix()

	var fresh []Task

	for _, t := range m.tasks {

		// keep if not expired OR not done
		if now-t.CreatedAt < 86400 || t.State != Done {
			fresh = append(fresh, t)
		}
	}

	m.tasks = fresh
	m.save()
}

func (m *TaskManager) startCleanupLoop() {

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			m.cleanupOldTasks()
		}
	}()
}
