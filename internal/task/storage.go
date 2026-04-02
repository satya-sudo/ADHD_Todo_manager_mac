package task

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	"focusbar/internal/logx"

	_ "github.com/mattn/go-sqlite3"
)

func (m *Manager) Load() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.ensureDBLocked(); err != nil {
		logx.Errorf("task load failed to open db path=%s err=%v", m.dbPath, err)
		m.tasks = []Task{}
		m.activeTaskID = ""
		return
	}

	tasks, err := m.loadTasksLocked()
	if err != nil {
		logx.Errorf("task load query failed path=%s err=%v", m.dbPath, err)
		m.tasks = []Task{}
		m.activeTaskID = ""
		return
	}

	if len(tasks) == 0 {
		migrated, migrateErr := m.migrateLegacyJSONLocked()
		if migrateErr != nil {
			logx.Errorf("task migration failed legacy=%s err=%v", m.legacyPath, migrateErr)
		} else if migrated {
			tasks, err = m.loadTasksLocked()
			if err != nil {
				logx.Errorf("task load after migration failed path=%s err=%v", m.dbPath, err)
				m.tasks = []Task{}
				m.activeTaskID = ""
				return
			}
		}
	}

	m.tasks = tasks
	logx.Infof("task load success path=%s count=%d", m.dbPath, len(m.tasks))

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
	if err := m.ensureDBLocked(); err != nil {
		logx.Errorf("task save open db failed path=%s err=%v", m.dbPath, err)
		return
	}

	tx, err := m.db.Begin()
	if err != nil {
		logx.Errorf("task save begin failed path=%s err=%v", m.dbPath, err)
		return
	}

	if _, err := tx.Exec(`DELETE FROM tasks`); err != nil {
		_ = tx.Rollback()
		logx.Errorf("task save clear failed path=%s err=%v", m.dbPath, err)
		return
	}

	stmt, err := tx.Prepare(`INSERT INTO tasks (id, title, state, created_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		logx.Errorf("task save prepare failed path=%s err=%v", m.dbPath, err)
		return
	}
	defer stmt.Close()

	for _, current := range m.tasks {
		if _, err := stmt.Exec(current.ID, current.Title, string(current.State), current.CreatedAt); err != nil {
			_ = tx.Rollback()
			logx.Errorf("task save insert failed path=%s task=%s err=%v", m.dbPath, current.ID, err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logx.Errorf("task save commit failed path=%s err=%v", m.dbPath, err)
		return
	}

	logx.Infof("task save success path=%s count=%d", m.dbPath, len(m.tasks))
}

func (m *Manager) ensureDBLocked() error {
	if m.db != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(m.dbPath), 0o755); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", m.dbPath)
	if err != nil {
		return err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			state TEXT NOT NULL,
			created_at INTEGER NOT NULL
		)
	`); err != nil {
		_ = db.Close()
		return err
	}

	m.db = db
	return nil
}

func (m *Manager) loadTasksLocked() ([]Task, error) {
	rows, err := m.db.Query(`SELECT id, title, state, created_at FROM tasks ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var current Task
		var state string
		if err := rows.Scan(&current.ID, &current.Title, &state, &current.CreatedAt); err != nil {
			return nil, err
		}
		current.State = State(state)
		tasks = append(tasks, current)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *Manager) migrateLegacyJSONLocked() (bool, error) {
	if m.legacyPath == "" {
		return false, nil
	}

	data, err := os.ReadFile(m.legacyPath)
	if err != nil {
		if os.IsNotExist(err) {
			logx.Infof("task legacy json not found path=%s", m.legacyPath)
			return false, nil
		}
		return false, err
	}

	var legacyTasks []Task
	if err := json.Unmarshal(data, &legacyTasks); err != nil {
		return false, err
	}

	if len(legacyTasks) == 0 {
		logx.Infof("task legacy json empty path=%s", m.legacyPath)
		return false, nil
	}

	m.tasks = legacyTasks
	m.saveLocked()
	logx.Infof("task legacy migration success legacy=%s count=%d", m.legacyPath, len(legacyTasks))
	return true, nil
}
