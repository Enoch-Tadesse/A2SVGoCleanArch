package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// CollectionTask is the name of the MongoDB collection for tasks
	CollectionTask = "tasks"
	// CollectionUser is the name of the MongoDB collection for users
	CollectionUser = "users"
)

// Task represents a task entity in the system
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`                     // Unique identifier for the task
	Title       string             `bson:"title" binding:"required" json:"title"`       // Title of the task
	Description string             `bson:"description" json:"description"`              // Detailed description of the task
	DueDate     time.Time          `bson:"due_date" binding:"required" json:"due_date"` // Due date for the task
	Status      string             `bson:"status" json:"status"`                        // Status of the task (e.g., "pending", "completed")
}

// TaskRepository defines the interface for interacting with the task persistence layer
type TaskRepository interface {
	// Create inserts a new task into the data store
	Create(c context.Context, task *Task) error
	// FetchByTaskID retrieves a task by its unique ID
	FetchByTaskID(c context.Context, taskID primitive.ObjectID) (Task, error)
	// FetchAllTasks retrieves all tasks from the data store
	FetchAllTasks(c context.Context) ([]Task, error)
	// DeleteByTaskID removes a task by its ID, returning the number of documents deleted
	DeleteByTaskID(c context.Context, taskID primitive.ObjectID) (int, error)
	// UpdateByTaskID updates an existing task, returning matched and modified counts
	UpdateByTaskID(c context.Context, task *Task) (int, int, error)
}

// TaskUsecase defines the business logic layer for task-related operations
type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByTaskID(c context.Context, taskID primitive.ObjectID) (Task, error)
	FetchAllTasks(c context.Context) ([]Task, error)
	DeleteByTaskID(c context.Context, taskID primitive.ObjectID) (int, error)
	UpdateByTaskID(c context.Context, task *Task) (int, int, error)
}

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"` // Unique identifier for the user
	Username string             `json:"username" bson:"username"` // Username of the user
	Password string             `json:"-" bson:"password"`        // Hashed password (excluded from JSON responses)
	IsAdmin  bool               `json:"is_admin" bson:"is_admin"` // Flag indicating if the user is an admin
}

// UserRepository defines the interface for interacting with the user persistence layer
type UserRepository interface {
	// Create inserts a new user into the data store
	Create(c context.Context, user *User) error
	// FetchByUserID retrieves a user by their unique ID
	FetchByUserID(c context.Context, userID primitive.ObjectID) (User, error)
	// FetchByUsername retrieves a user by their username
	FetchByUsername(c context.Context, username string) (User, error)
	// FetchAllUsers retrieves all users from the data store
	FetchAllUsers(c context.Context) ([]User, error)
	// PromoteByUserID sets the IsAdmin flag to true for the specified user
	PromoteByUserID(c context.Context, userID primitive.ObjectID) (int, error)
}

// UserUsercase defines the business logic layer for user-related operations
type UserUsercase interface {
	Create(c context.Context, user *User) error
	FetchByUserID(c context.Context, userID primitive.ObjectID) (User, error)
	FetchByUsername(c context.Context, username string) (User, error)
	FetchAllUsers(c context.Context) ([]User, error)
	PromoteByUserID(c context.Context, userID primitive.ObjectID) (int, error)
}
