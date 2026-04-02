package app

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"focusbar/internal/logx"
)

var ErrAlreadyRunning = errors.New("focusbar is already running")

type InstanceLock struct {
	file *os.File
}

func AcquireInstanceLock() (*InstanceLock, error) {
	path := defaultLockPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, ErrAlreadyRunning
		}
		return nil, err
	}

	logx.Infof("instance lock acquired path=%s", path)
	return &InstanceLock{file: file}, nil
}

func (l *InstanceLock) Release() error {
	if l == nil || l.file == nil {
		return nil
	}

	path := l.file.Name()
	if err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN); err != nil {
		_ = l.file.Close()
		return err
	}

	if err := l.file.Close(); err != nil {
		return err
	}

	logx.Infof("instance lock released path=%s", path)
	l.file = nil
	return nil
}

func defaultLockPath() string {
	if path := os.Getenv("FOCUSBAR_LOCK_PATH"); path != "" {
		return path
	}

	return filepath.Join(defaultDataDir(), "focusbar.lock")
}
