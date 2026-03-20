package main

import (
	"focusbar/internal/app"

	"github.com/getlantern/systray"
)

var application = app.New()

func main() {
	systray.Run(application.OnReady, application.OnExit)
}
