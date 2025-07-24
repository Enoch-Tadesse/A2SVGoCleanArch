package usecases

import (
	"context"
	"strings"
	"time"

	domain "github.com/A2SVTask7/Domain"
)

// taskUsecase implements the domain.TaskUsecase interface
type taskUsecase struct {
	taskRepository domain.TaskRepository // Repository for task data operations
	contextTimeout time.Duration         // Timeout duration for each usecase operation
}

// NewTaskUsecase creates a new instance of taskUsecase
func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

// Create adds a new task using the repository with a timeout context
func (tu *taskUsecase) Create(c context.Context, task *domain.Task) error {

	if task.DueDate.Before(time.Now()) {
		return domain.ErrInvalidDueDate
	}
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.Create(ctx, task)
}

// UpdateByTaskID updates an existing task using its ID
// Returns the number of matched and modified documents
func (tu *taskUsecase) UpdateByTaskID(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// Normalize status field
	task.Status = strings.ToLower(strings.TrimSpace(task.Status))

	// Validate due date
	if task.DueDate.Before(time.Now()) {
		return domain.ErrInvalidDueDate
	}

	matched, modified, err := tu.taskRepository.UpdateByTaskID(ctx, task)
	if err != nil {
		return err
	}

	if matched == 0 {
		return domain.ErrTaskNotFound
	}

	if modified == 0 {
		return domain.ErrNoChangesMade
	}

	return nil
}

// DeleteByTaskID deletes a task using its ID
// Returns the number of documents deleted
func (tu *taskUsecase) DeleteByTaskID(c context.Context, taskID string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	count, err := tu.taskRepository.DeleteByTaskID(ctx, taskID)
	if count == 0 {
		return domain.ErrTaskNotFound
	}
	return err
}

// FetchByTaskID retrieves a single task by its ID
func (tu *taskUsecase) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.FetchByTaskID(ctx, taskID)
}

// FetchAllTasks retrieves all tasks from the repository
func (tu *taskUsecase) FetchAllTasks(c context.Context) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.FetchAllTasks(ctx)
}
