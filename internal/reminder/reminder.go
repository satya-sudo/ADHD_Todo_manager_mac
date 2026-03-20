package reminder

import "os/exec"

func Notify(task string) {
	_ = exec.Command(
		"osascript",
		"-e",
		`display notification "Still working on: `+task+`" with title "⚡ Focus Check"`,
	).Run()
}
