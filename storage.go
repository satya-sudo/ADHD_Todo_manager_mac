package main

import (
	"encoding/json"
	"os"
)

const taskFile = "tasks.json"

func (m *TaskManager) load() {

	data, err := os.ReadFile(taskFile)
	if err != nil {
		// file may not exist yet → fine
		m.tasks = []Task{}
		return
	}

	json.Unmarshal(data, &m.tasks)

	// IMPORTANT: restore activeTaskID
	for _, t := range m.tasks {
		if t.State == Working {
			m.activeTaskID = t.ID
			return
		}
	}
}
func (m *TaskManager) save() {

	data, err := json.MarshalIndent(m.tasks, "", " ")
	if err != nil {
		return
	}

	os.WriteFile(taskFile, data, 0644)
}
