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

	application := app.New()
	systray.Run(application.OnReady, application.OnExit)
}
