package main

import (
	"os/exec"
)

//func startReminder() {
//
//	go func() {
//
//		for {
//
//			time.Sleep(2 * time.Hour)
//
//			active := getActiveTask()
//
//			if active != "" {
//
//				notify(active)
//
//			}
//
//		}
//
//	}()
//
//}

func notify(task string) {

	exec.Command(
		"osascript",
		"-e",
		`display notification "Still working on: `+task+`" with title "⚡ Focus Check"`,
	).Run()

}
