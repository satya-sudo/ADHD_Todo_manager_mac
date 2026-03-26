package task

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSQLiteRoundTrip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "focusbar.db")

	manager := NewManager(dbPath, "", nil, nil, nil)
	manager.tasks = []Task{
		{ID: "1", Title: "first", State: Todo, CreatedAt: 100},
		{ID: "2", Title: "second", State: Working, CreatedAt: 200},
	}
	manager.activeTaskID = "2"
	manager.saveLocked()
	t.Cleanup(func() { _ = manager.Close() })

	loaded := NewManager(dbPath, "", nil, nil, nil)
	loaded.Load()
	t.Cleanup(func() { _ = loaded.Close() })

	if len(loaded.tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(loaded.tasks))
	}
	if loaded.activeTaskID != "2" {
		t.Fatalf("expected active task 2, got %q", loaded.activeTaskID)
	}
	if loaded.tasks[1].Title != "second" {
		t.Fatalf("unexpected second task title %q", loaded.tasks[1].Title)
	}
}

func TestSQLiteMigratesLegacyJSON(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "focusbar.db")
	legacyPath := filepath.Join(tmpDir, "tasks.json")
	legacyTasks := []Task{
		{ID: "a", Title: "legacy todo", State: Todo, CreatedAt: 10},
		{ID: "b", Title: "legacy active", State: Working, CreatedAt: 20},
	}

	data, err := json.Marshal(legacyTasks)
	if err != nil {
		t.Fatalf("marshal legacy tasks: %v", err)
	}
	if err := osWriteFile(legacyPath, data); err != nil {
		t.Fatalf("write legacy file: %v", err)
	}

	manager := NewManager(dbPath, legacyPath, nil, nil, nil)
	manager.Load()
	t.Cleanup(func() { _ = manager.Close() })

	if len(manager.tasks) != 2 {
		t.Fatalf("expected 2 migrated tasks, got %d", len(manager.tasks))
	}
	if manager.activeTaskID != "b" {
		t.Fatalf("expected active task b after migration, got %q", manager.activeTaskID)
	}

	reloaded := NewManager(dbPath, "", nil, nil, nil)
	reloaded.Load()
	t.Cleanup(func() { _ = reloaded.Close() })

	if len(reloaded.tasks) != 2 {
		t.Fatalf("expected 2 persisted tasks after reload, got %d", len(reloaded.tasks))
	}
}

func osWriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}
