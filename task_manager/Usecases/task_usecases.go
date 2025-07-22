package usecases

import (
	"context"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.Create(ctx, task)
}

// UpdateByTaskID updates an existing task using its ID
// Returns the number of matched and modified documents
func (tu *taskUsecase) UpdateByTaskID(c context.Context, task *domain.Task) (int, int, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.UpdateByTaskID(ctx, task)
}

// DeleteByTaskID deletes a task using its ID
// Returns the number of documents deleted
func (tu *taskUsecase) DeleteByTaskID(c context.Context, taskID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.DeleteByTaskID(ctx, taskID)
}

// FetchByTaskID retrieves a single task by its ID
func (tu *taskUsecase) FetchByTaskID(c context.Context, taskID primitive.ObjectID) (domain.Task, error) {
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
