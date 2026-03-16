package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

var focusStart time.Time

func startFocusTimer(title string) {

	focusStart = time.Now()

	go func(task string) {

		for {

			elapsed := time.Since(focusStart)

			minutes := int(elapsed.Minutes())

			systray.SetTitle(fmt.Sprintf("⚡ %dm %s", minutes, task))

			time.Sleep(1 * time.Minute)

		}

	}(title)

}
