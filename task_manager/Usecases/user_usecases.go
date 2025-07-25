package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/A2SVTask7/Domain"
	infrastructure "github.com/A2SVTask7/Infrastructure"
)

// userUsecase implements the domain.UserUsercase interface
type userUsecase struct {
	userRepository domain.UserRepository // Repository for user data operations
	contextTimeout time.Duration         // Timeout duration for usecase operations
}

// NewUserUsecase creates a new instance of userUsecase
func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

// Login is a usecase to help compare password, validate user and generate claims
func (uu *userUsecase) Login(ctx context.Context, username, password string) (domain.User, string, error) {
	user, err := uu.userRepository.FetchByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, "", domain.ErrUserNotFound
		}
		return domain.User{}, "", err
	}

	if err := infrastructure.ComparePassword(user.Password, password); err != nil {
		return domain.User{}, "", domain.ErrIncorrectPassword // define this in domain if you want
	}

	claims := map[string]any{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	js := infrastructure.NewJWTService(infrastructure.AppConfig.JWTSecret)
	token, err := js.Generate(claims)
	if err != nil {
		return domain.User{}, "", err
	}

	return user, token, nil
}

// CheckIfUsernameExists checks if user exist or not
func (uu *userUsecase) CheckIfUsernameExists(c context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.CheckIfUsernameExists(ctx, username)
}

// CountUsers counts the number of users in database
func (uu *userUsecase) CountUsers(c context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.CountUsers(ctx)
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

	hashedString, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedString

	exists, err := uu.CheckIfUsernameExists(c, user.Username)
	if err != nil {
		return errors.New("failed to check if username exists")
	}
	if exists {
		return domain.ErrUserAlreadyExists
	}

	count, err := uu.CountUsers(c)
	if err != nil {
		return errors.New("failed to count users")
	}

	user.IsAdmin = count == 0
	return uu.userRepository.Create(ctx, user)
}

// PromoteByUserID promotes a user to admin by setting IsAdmin to true
// Returns the number of documents modified
func (uu *userUsecase) PromoteByUserID(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	count, err := uu.userRepository.PromoteByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err.Error())
	}
	if count == 0 {
		return domain.ErrUserNotFound
	}
	return nil

}

// FetchAllUsers retrieves all users from the repository
func (uu *userUsecase) FetchAllUsers(c context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.FetchAllUsers(ctx)
}

// FetchByUserID retrieves a user by their unique ID
func (uu *userUsecase) FetchByUserID(c context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.FetchByUserID(ctx, userID)
}
