package domain

import "errors"

// custom errors
var (
	ErrInvalidTaskID  = errors.New("invalid task id")
	ErrInvalidDueDate = errors.New("due date cannot be in the past")
	ErrTaskNotFound   = errors.New("task not found")
	ErrNoChangesMade  = errors.New("no changes were made")
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrIncorrectPassword = errors.New("incorrect password")
)
