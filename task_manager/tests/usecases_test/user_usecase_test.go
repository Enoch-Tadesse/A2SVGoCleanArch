package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	domain "github.com/A2SVTask7/Domain"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	usecases "github.com/A2SVTask7/Usecases"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	mockRepo    *MockUserRepository
	userUsecase domain.UserUsecase
	ctx         context.Context
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.userUsecase = usecases.NewUserUsecase(s.mockRepo, 2*time.Second)
	s.ctx = context.Background()
}

var sampleUser = domain.User{
	ID:       "user-id-123",
	Username: "testuser",
	Password: "$2a$10$somehashedvalue",
	IsAdmin:  false,
}

func (s *UserUsecaseTestSuite) TestLogin_Success() {
	user := sampleUser
	hashed, _ := infrastructure.HashPassword("password")
	user.Password = hashed

	s.mockRepo.On("FetchByUsername", mock.Anything, "testuser").Return(user, nil)

	u, token, err := s.userUsecase.Login(s.ctx, "testuser", "password")

	s.NoError(err)
	s.NotEmpty(token)
	s.Equal(user.ID, u.ID)
}

func (s *UserUsecaseTestSuite) TestLogin_WrongPassword() {
	s.mockRepo.On("FetchByUsername", mock.Anything, "testuser").Return(sampleUser, nil)

	_, _, err := s.userUsecase.Login(s.ctx, "testuser", "wrongpass")
	s.EqualError(err, domain.ErrIncorrectPassword.Error())
}

func (s *UserUsecaseTestSuite) TestLogin_UserNotFound() {
	s.mockRepo.On("FetchByUsername", mock.Anything, "unknown").Return(domain.User{}, domain.ErrUserNotFound)

	_, _, err := s.userUsecase.Login(s.ctx, "unknown", "pass")
	s.EqualError(err, domain.ErrUserNotFound.Error())
}

func (s *UserUsecaseTestSuite) TestCreate_Success() {
	user := domain.User{Username: "newuser", Password: "securepass"}

	s.mockRepo.On("CheckIfUsernameExists", mock.Anything, "newuser").Return(false, nil)
	s.mockRepo.On("CountUsers", mock.Anything).Return(0, nil)
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := s.userUsecase.Create(s.ctx, &user)

	s.NoError(err)
	s.True(user.IsAdmin)
}

func (s *UserUsecaseTestSuite) TestCreate_UsernameExists() {
	user := domain.User{Username: "existing", Password: "pass"}

	s.mockRepo.On("CheckIfUsernameExists", mock.Anything, "existing").Return(true, nil)

	err := s.userUsecase.Create(s.ctx, &user)
	s.EqualError(err, domain.ErrUserAlreadyExists.Error())
}

func (s *UserUsecaseTestSuite) TestCountUsers() {
	s.mockRepo.On("CountUsers", mock.Anything).Return(5, nil)

	count, err := s.userUsecase.CountUsers(s.ctx)
	s.NoError(err)
	s.Equal(5, count)
}

func (s *UserUsecaseTestSuite) TestCheckIfUsernameExists() {
	s.mockRepo.On("CheckIfUsernameExists", mock.Anything, "checkuser").Return(true, nil)

	exists, err := s.userUsecase.CheckIfUsernameExists(s.ctx, "checkuser")
	s.NoError(err)
	s.True(exists)
}

func (s *UserUsecaseTestSuite) TestPromoteByUserID_Success() {
	s.mockRepo.On("PromoteByUserID", mock.Anything, "user-id-123").Return(1, nil)

	err := s.userUsecase.PromoteByUserID(s.ctx, "user-id-123")
	s.NoError(err)
}

func (s *UserUsecaseTestSuite) TestPromoteByUserID_NotFound() {
	s.mockRepo.On("PromoteByUserID", mock.Anything, "nonexistent").Return(0, nil)

	err := s.userUsecase.PromoteByUserID(s.ctx, "nonexistent")
	s.EqualError(err, domain.ErrUserNotFound.Error())
}

func (s *UserUsecaseTestSuite) TestFetchAllUsers() {
	s.mockRepo.On("FetchAllUsers", mock.Anything).Return([]domain.User{sampleUser}, nil)

	users, err := s.userUsecase.FetchAllUsers(s.ctx)
	s.NoError(err)
	s.Len(users, 1)
	s.Equal(sampleUser.ID, users[0].ID)
}

func (s *UserUsecaseTestSuite) TestFetchByUsername() {
	s.mockRepo.On("FetchByUsername", mock.Anything, "testuser").Return(sampleUser, nil)

	user, err := s.userUsecase.FetchByUsername(s.ctx, "testuser")
	s.NoError(err)
	s.Equal("testuser", user.Username)
}

func (s *UserUsecaseTestSuite) TestFetchByUserID() {
	s.mockRepo.On("FetchByUserID", mock.Anything, "user-id-123").Return(sampleUser, nil)

	user, err := s.userUsecase.FetchByUserID(s.ctx, "user-id-123")
	s.NoError(err)
	s.Equal("user-id-123", user.ID)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
