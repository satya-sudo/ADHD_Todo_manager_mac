package main

import (
	"focusbar/internal/app"
	"focusbar/internal/logx"

	"github.com/getlantern/systray"
)

func main() {
	if err := logx.Init(); err != nil {
		panic(err)
	}

	logx.Infof("starting focusbar")
	logx.Infof("log file path=%s", logx.Path())

	lock, err := app.AcquireInstanceLock()
	if err != nil {
		if err == app.ErrAlreadyRunning {
			logx.Infof("focusbar launch skipped: another instance is already running")
			return
		}
		panic(err)
	}
	defer func() {
		if err := lock.Release(); err != nil {
			logx.Errorf("instance lock release failed err=%v", err)
		}
	}()

	application := app.New()
	systray.Run(application.OnReady, application.OnExit)
}
