package model

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskStatusRequest struct {
	Status TaskStatus `json:"status" binding:"required"`
}

type TaskStats struct {
	Total      int `json:"total"`
	Todo       int `json:"todo"`
	Done       int `json:"done"`
	InProgress int `json:"in_progress"`
}

func (status TaskStatus) IsValid() bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
