package tray

import "unicode/utf8"

const maxTitleRunes = 18

func IdleTitle() string {
	return "⚡ Ready to focus"
}

func NudgeTitle() string {
	return "⚡ Pick something"
}

func ResumeTitle() string {
	return "⚡ Resume one"
}

func OverloadTitle() string {
	return "⚡ Pick just one"
}

func TaskTitle(title string) string {
	return "⚡ " + truncate(title, maxTitleRunes)
}

func TimerTitle(timerText string, title string) string {
	return timerText + " " + truncate(title, 12)
}

func truncate(value string, maxRunes int) string {
	if utf8.RuneCountInString(value) <= maxRunes {
		return value
	}

	runes := []rune(value)
	if maxRunes <= 1 {
		return string(runes[:maxRunes])
	}

	return string(runes[:maxRunes-1]) + "…"
}
