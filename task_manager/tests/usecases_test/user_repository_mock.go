package usecases_test

import (
	"context"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(c context.Context, user *domain.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}

func (m *MockUserRepository) FetchByUserID(c context.Context, userID string) (domain.User, error) {
	args := m.Called(c, userID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) FetchByUsername(c context.Context, username string) (domain.User, error) {
	args := m.Called(c, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) FetchAllUsers(c context.Context) ([]domain.User, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) PromoteByUserID(c context.Context, userID string) (int, error) {
	args := m.Called(c, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) CountUsers(c context.Context) (int, error) {
	args := m.Called(c)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) CheckIfUsernameExists(c context.Context, username string) (bool, error) {
	args := m.Called(c, username)
	return args.Bool(0), args.Error(1)
}
