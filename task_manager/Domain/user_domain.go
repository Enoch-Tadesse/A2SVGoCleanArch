package domain

import (
	"context"
)

// User represents a user in the system
type User struct {
	ID       string // Unique identifier for the user
	Username string // Username of the user
	Password string // Hashed password (excluded from JSON responses)
	IsAdmin  bool   // Flag indicating if the user is an admin
}

// UserRepository defines the interface for interacting with the user persistence layer
type UserRepository interface {
	// Create inserts a new user into the data store
	Create(c context.Context, user *User) error
	// FetchByUserID retrieves a user by their unique ID
	FetchByUserID(c context.Context, userID string) (User, error)
	// FetchByUsername retrieves a user by their username
	FetchByUsername(c context.Context, username string) (User, error)
	// FetchAllUsers retrieves all users from the data store
	FetchAllUsers(c context.Context) ([]User, error)
	// PromoteByUserID sets the IsAdmin flag to true for the specified user
	PromoteByUserID(c context.Context, userID string) (int, error)
}

// UserUsercase defines the business logic layer for user-related operations
type UserUsercase interface {
	Create(c context.Context, user *User) error
	FetchByUserID(c context.Context, userID string) (User, error)
	FetchByUsername(c context.Context, username string) (User, error)
	FetchAllUsers(c context.Context) ([]User, error)
	PromoteByUserID(c context.Context, userID string) (int, error)
}
