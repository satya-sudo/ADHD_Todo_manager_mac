package logx

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	mu       sync.RWMutex
	logger   = log.New(io.Discard, "", log.LstdFlags|log.Lmicroseconds)
	logPath  string
	logFile  *os.File
	initOnce sync.Once
)

func Init() error {
	var initErr error

	initOnce.Do(func() {
		path := os.Getenv("FOCUSBAR_LOG_PATH")
		if path == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				initErr = err
				return
			}

			path = filepath.Join(home, "Library", "Logs", "Focusbar", "focusbar.log")
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			initErr = err
			return
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			initErr = err
			return
		}

		mu.Lock()
		defer mu.Unlock()

		logPath = path
		logFile = file
		logger.SetOutput(file)
		logger.Printf("logger initialized path=%s", path)
	})

	return initErr
}

func Path() string {
	mu.RLock()
	defer mu.RUnlock()

	return logPath
}

func Infof(format string, args ...any) {
	mu.RLock()
	defer mu.RUnlock()

	logger.Printf("INFO "+format, args...)
}

func Errorf(format string, args ...any) {
	mu.RLock()
	defer mu.RUnlock()

	logger.Printf("ERROR "+format, args...)
}
