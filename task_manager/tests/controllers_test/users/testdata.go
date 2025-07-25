package users

// this file contails dommon shared data and struct within users

import (
	domain "github.com/A2SVTask7/Domain"
)

// request struct
type UserRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserListTestCase struct {
	Name          string              // name of the test
	MockSetup     func()              // mock strup
	Expected      int                 // expected status
	ValidateSlice func([]domain.User) // validation datas of a slice
	Validate      func(domain.User)   // validation of a single data
	Payload       UserRequest         // payload if it is a post request
}

// User represents a user in the system
type User struct {
	ID       string // Unique identifier for the user
	Username string // Username of the user
	Password string // Hashed password (excluded from JSON responses)
	IsAdmin  bool   // Flag indicating if the user is an admin
}

// valid sample users
var sampleUsers = []domain.User{
	{
		ID:       "user1",
		Username: "alice",
		Password: "$2a$10$hashedpassword1", // bcrypt-style fake hash
		IsAdmin:  false,
	},
	{
		ID:       "user2",
		Username: "bob",
		Password: "$2a$10$hashedpassword2",
		IsAdmin:  true,
	},
	{
		ID:       "user3",
		Username: "charlie",
		Password: "$2a$10$hashedpassword3",
		IsAdmin:  false,
	},
}
