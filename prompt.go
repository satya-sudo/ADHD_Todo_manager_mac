package main

import (
	"os/exec"
	"strings"
)

func promptTask() string {

	cmd := exec.Command(
		"osascript",
		"-e",
		`text returned of (display dialog "⚡ Add Task" default answer "")`,
	)

	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
