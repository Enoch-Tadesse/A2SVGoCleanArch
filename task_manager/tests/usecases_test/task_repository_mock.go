package usecases_test

import (
	"context"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository is a mock implementation of the TaskRepository interface
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(c context.Context, task *domain.Task) error {
	args := m.Called(c, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	args := m.Called(c, taskID)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FetchAllTasks(c context.Context) ([]domain.Task, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteByTaskID(c context.Context, taskID string) (int, error) {
	args := m.Called(c, taskID)
	return args.Int(0), args.Error(1)
}

func (m *MockTaskRepository) UpdateByTaskID(c context.Context, task *domain.Task) (int, int, error) {
	args := m.Called(c, task)
	return args.Int(0), args.Int(1), args.Error(2)
}
