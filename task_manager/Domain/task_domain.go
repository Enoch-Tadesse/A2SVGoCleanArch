package domain

import (
	"context"
	"time"
)

// Task represents a task entity in the system
type Task struct {
	ID          string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

// TaskRepository defines the interface for interacting with the task persistence layer
type TaskRepository interface {
	// Create inserts a new task into the data store
	Create(c context.Context, task *Task) error
	// FetchByTaskID retrieves a task by its unique ID
	FetchByTaskID(c context.Context, taskID string) (Task, error)
	// FetchAllTasks retrieves all tasks from the data store
	FetchAllTasks(c context.Context) ([]Task, error)
	// DeleteByTaskID removes a task by its ID, returning the number of documents deleted
	DeleteByTaskID(c context.Context, taskID string) (int, error)
	// UpdateByTaskID updates an existing task, returning matched and modified counts
	UpdateByTaskID(c context.Context, task *Task) (int, int, error)
}

// TaskUsecase defines the business logic layer for task-related operations
type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByTaskID(c context.Context, taskID string) (Task, error)
	FetchAllTasks(c context.Context) ([]Task, error)
	DeleteByTaskID(c context.Context, taskID string) error
	UpdateByTaskID(c context.Context, task *Task) error
}
