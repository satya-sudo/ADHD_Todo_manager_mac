package main

import (
	"encoding/json"
	"os"
)

var tasks []Task

func loadTasks() {

	file, err := os.ReadFile("tasks_today.json")
	if err != nil {
		return
	}

	json.Unmarshal(file, &tasks)
}

func saveTasks() {

	data, _ := json.MarshalIndent(tasks, "", " ")

	os.WriteFile("tasks_today.json", data, 0644)
}
