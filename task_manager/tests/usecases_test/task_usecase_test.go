package usecases_test

import (
	"context"
	"testing"
	"time"

	domain "github.com/A2SVTask7/Domain"
	usecases "github.com/A2SVTask7/Usecases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
	mockRepo    *MockTaskRepository
	taskUsecase domain.TaskUsecase
	ctx         context.Context
}

func (s *TaskUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockTaskRepository)
	s.taskUsecase = usecases.NewTaskUsecase(s.mockRepo, time.Second*2)
	s.ctx = context.Background()
}

var sampleTask = domain.Task{
	ID:          "task-id-123",
	Title:       "Test Task",
	Description: "Test Description",
	DueDate:     time.Now().Add(time.Hour),
	Status:      "Pending",
}

func (s *TaskUsecaseTestSuite) TestCreate_Success() {
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil)

	err := s.taskUsecase.Create(s.ctx, &sampleTask)
	s.NoError(err)
	s.mockRepo.AssertCalled(s.T(), "Create", mock.Anything, &sampleTask)
}

func (s *TaskUsecaseTestSuite) TestCreate_InvalidDueDate() {
	task := sampleTask
	task.DueDate = time.Now().Add(-time.Hour)

	err := s.taskUsecase.Create(s.ctx, &task)
	s.EqualError(err, domain.ErrInvalidDueDate.Error())
	s.mockRepo.AssertNotCalled(s.T(), "Create")
}

func (s *TaskUsecaseTestSuite) TestUpdateByTaskID_Success() {
	task := sampleTask
	s.mockRepo.On("UpdateByTaskID", mock.Anything, &task).Return(1, 1, nil)

	err := s.taskUsecase.UpdateByTaskID(s.ctx, &task)
	s.NoError(err)
	s.mockRepo.AssertCalled(s.T(), "UpdateByTaskID", mock.Anything, &task)
}

func (s *TaskUsecaseTestSuite) TestUpdateByTaskID_NoMatch() {
	task := sampleTask
	s.mockRepo.On("UpdateByTaskID", mock.Anything, &task).Return(0, 0, nil)

	err := s.taskUsecase.UpdateByTaskID(s.ctx, &task)
	s.EqualError(err, domain.ErrTaskNotFound.Error())
}

func (s *TaskUsecaseTestSuite) TestUpdateByTaskID_NoChange() {
	task := sampleTask
	s.mockRepo.On("UpdateByTaskID", mock.Anything, &task).Return(1, 0, nil)

	err := s.taskUsecase.UpdateByTaskID(s.ctx, &task)
	s.EqualError(err, domain.ErrNoChangesMade.Error())
}

func (s *TaskUsecaseTestSuite) TestDeleteByTaskID_Success() {
	s.mockRepo.On("DeleteByTaskID", mock.Anything, "task-id-123").Return(1, nil)

	err := s.taskUsecase.DeleteByTaskID(s.ctx, "task-id-123")
	s.NoError(err)
}

func (s *TaskUsecaseTestSuite) TestDeleteByTaskID_NotFound() {
	s.mockRepo.On("DeleteByTaskID", mock.Anything, "non-existent-id").Return(0, nil)

	err := s.taskUsecase.DeleteByTaskID(s.ctx, "non-existent-id")
	s.EqualError(err, domain.ErrTaskNotFound.Error())
}

func (s *TaskUsecaseTestSuite) TestFetchByTaskID_Success() {
	s.mockRepo.On("FetchByTaskID", mock.Anything, "task-id-123").Return(sampleTask, nil)

	task, err := s.taskUsecase.FetchByTaskID(s.ctx, "task-id-123")
	s.NoError(err)
	s.Equal(sampleTask.ID, task.ID)
}

func (s *TaskUsecaseTestSuite) TestFetchAllTasks_Success() {
	s.mockRepo.On("FetchAllTasks", mock.Anything).Return([]domain.Task{sampleTask}, nil)

	tasks, err := s.taskUsecase.FetchAllTasks(s.ctx)
	s.NoError(err)
	s.Len(tasks, 1)
	s.Equal(sampleTask.ID, tasks[0].ID)
}

func TestTaskUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}
