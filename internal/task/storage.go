package task

import (
	"encoding/json"
	"os"
	"path/filepath"

	"focusbar/internal/logx"
)

func (m *Manager) Load() {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.storagePath)
	if err != nil {
		logx.Infof("task load no existing storage path=%s err=%v", m.storagePath, err)
		m.tasks = []Task{}
		m.activeTaskID = ""
		return
	}

	if err := json.Unmarshal(data, &m.tasks); err != nil {
		logx.Errorf("task load failed path=%s err=%v", m.storagePath, err)
		m.tasks = []Task{}
		m.activeTaskID = ""
		return
	}

	logx.Infof("task load success path=%s count=%d", m.storagePath, len(m.tasks))

	m.activeTaskID = ""
	for _, current := range m.tasks {
		if current.State == Working {
			m.activeTaskID = current.ID
			return
		}
	}
}

func (m *Manager) save() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.saveLocked()
}

func (m *Manager) saveLocked() {
	if err := os.MkdirAll(filepath.Dir(m.storagePath), 0o755); err != nil {
		logx.Errorf("task save mkdir failed path=%s err=%v", m.storagePath, err)
		return
	}

	data, err := json.MarshalIndent(m.tasks, "", " ")
	if err != nil {
		logx.Errorf("task save marshal failed err=%v", err)
		return
	}

	if err := os.WriteFile(m.storagePath, data, 0o644); err != nil {
		logx.Errorf("task save write failed path=%s err=%v", m.storagePath, err)
		return
	}

	logx.Infof("task save success path=%s count=%d", m.storagePath, len(m.tasks))
}
