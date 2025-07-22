package usecases

import (
	"context"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// userUsecase implements the domain.UserUsercase interface
type userUsecase struct {
	userRepository domain.UserRepository // Repository for user data operations
	contextTimeout time.Duration         // Timeout duration for usecase operations
}

// NewUserUsecase creates a new instance of userUsecase
func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsercase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

// FetchByUsername retrieves a user by their username
func (uu *userUsecase) FetchByUsername(c context.Context, username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.FetchByUsername(ctx, username)
}

// Create registers a new user using the repository
func (uu *userUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.Create(ctx, user)
}

// PromoteByUserID promotes a user to admin by setting IsAdmin to true
// Returns the number of documents modified
func (uu *userUsecase) PromoteByUserID(c context.Context, userID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.PromoteByUserID(ctx, userID)
}

// FetchAllUsers retrieves all users from the repository
func (uu *userUsecase) FetchAllUsers(c context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.FetchAllUsers(ctx)
}

// FetchByUserID retrieves a user by their unique ID
func (uu *userUsecase) FetchByUserID(c context.Context, userID primitive.ObjectID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.FetchByUserID(ctx, userID)
}
