package task

import (
	"encoding/json"
	"os"
)

func (m *Manager) Load() {
	data, err := os.ReadFile(m.storagePath)
	if err != nil {
		m.tasks = []Task{}
		m.activeTaskID = ""
		return
	}

	if err := json.Unmarshal(data, &m.tasks); err != nil {
		m.tasks = []Task{}
		m.activeTaskID = ""
		return
	}

	m.activeTaskID = ""
	for _, current := range m.tasks {
		if current.State == Working {
			m.activeTaskID = current.ID
			return
		}
	}
}

func (m *Manager) save() {
	data, err := json.MarshalIndent(m.tasks, "", " ")
	if err != nil {
		return
	}

	_ = os.WriteFile(m.storagePath, data, 0o644)
}
