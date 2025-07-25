package mocks

import (
	"context"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Create(c context.Context, user *domain.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}

func (m *MockUserUsecase) FetchByUserID(c context.Context, userID string) (domain.User, error) {
	args := m.Called(c, userID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUsecase) FetchByUsername(c context.Context, username string) (domain.User, error) {
	args := m.Called(c, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUsecase) FetchAllUsers(c context.Context) ([]domain.User, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserUsecase) PromoteByUserID(c context.Context, userID string) error {
	args := m.Called(c, userID)
	return args.Error(0)
}

func (m *MockUserUsecase) CountUsers(c context.Context) (int, error) {
	args := m.Called(c)
	return args.Int(0), args.Error(1)
}

func (m *MockUserUsecase) CheckIfUsernameExists(c context.Context, username string) (bool, error) {
	args := m.Called(c, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserUsecase) Login(ctx context.Context, username, password string) (domain.User, string, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(domain.User), args.String(1), args.Error(2)
}
