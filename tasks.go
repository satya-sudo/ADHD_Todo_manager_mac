package main

import (
	"time"

	"github.com/google/uuid"
)

type TaskState string

const (
	Todo    TaskState = "todo"
	Working TaskState = "working"
	Paused  TaskState = "paused"
	Done    TaskState = "done"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	State     TaskState `json:"state"`
	CreatedAt int64     `json:"created_at"`
}

func NewTask(title string) Task {
	return Task{
		ID:        uuid.NewString(),
		Title:     title,
		State:     Todo,
		CreatedAt: time.Now().Unix(),
	}
}

func FindTaskByID(tasks []Task, id string) *Task {

	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}

func GetActiveTask(tasks []Task) *Task {

	for i := range tasks {
		if tasks[i].State == Working {
			return &tasks[i]
		}
	}
	return nil
}

func FilterTasksByState(tasks []Task, state TaskState) []Task {

	var result []Task

	for _, t := range tasks {
		if t.State == state {
			result = append(result, t)
		}
	}

	return result
}

func RemoveTaskByID(tasks []Task, id string) []Task {

	var result []Task

	for _, t := range tasks {
		if t.ID != id {
			result = append(result, t)
		}
	}

	return result
}
