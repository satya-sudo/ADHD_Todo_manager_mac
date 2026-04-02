package app

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestAcquireInstanceLockBlocksSecondAcquire(t *testing.T) {
	lockPath := filepath.Join(t.TempDir(), "focusbar.lock")
	t.Setenv("FOCUSBAR_LOCK_PATH", lockPath)

	first, err := AcquireInstanceLock()
	if err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	t.Cleanup(func() {
		if err := first.Release(); err != nil {
			t.Fatalf("release failed: %v", err)
		}
	})

	second, err := AcquireInstanceLock()
	if second != nil {
		t.Fatalf("expected second lock to fail, got lock=%v", second)
	}
	if !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("expected ErrAlreadyRunning, got %v", err)
	}
}
