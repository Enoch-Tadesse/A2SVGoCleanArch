package mocks

import (
	"context"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) Create(c context.Context, task *domain.Task) error {
	args := m.Called(c, task)
	return args.Error(0)
}

func (m *MockTaskUsecase) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	args := m.Called(c, taskID)
	return args.Get(0).(domain.Task), args.Error(1)
}
func (m *MockTaskUsecase) FetchAllTasks(c context.Context) ([]domain.Task, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.Task), args.Error(1)
}
func (m *MockTaskUsecase) DeleteByTaskID(c context.Context, taskID string) error {
	args := m.Called(c, taskID)
	return args.Error(0)
}
func (m *MockTaskUsecase) UpdateByTaskID(c context.Context, task *domain.Task) error {
	args := m.Called(c, task)
	return args.Error(0)
}
