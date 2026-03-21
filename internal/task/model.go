package task

import (
	"time"

	"github.com/google/uuid"
)

type State string
type TaskState = State

const (
	Todo    State = "todo"
	Working State = "working"
	Paused  State = "paused"
	Done    State = "done"
)

type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	State     State  `json:"state"`
	CreatedAt int64  `json:"created_at"`
}

func New(title string) Task {
	return Task{
		ID:        uuid.NewString(),
		Title:     title,
		State:     Todo,
		CreatedAt: time.Now().Unix(),
	}
}

func FindByID(tasks []Task, id string) *Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}

	return nil
}

func FilterByState(tasks []Task, state State) []Task {
	var filtered []Task

	for _, current := range tasks {
		if current.State == state {
			filtered = append(filtered, current)
		}
	}

	return filtered
}
