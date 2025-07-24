package controllers

import (
	"errors"
	"net/http"

	domain "github.com/A2SVTask7/Domain"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	"github.com/gin-gonic/gin"
)

// UserController handles user-related HTTP endpoints
type UserController struct {
	UserUsecase domain.UserUsecase
}

// GetUserByID handles GET /users/:id
// Fetches a user by their ObjectID
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	user, err := uc.UserUsecase.FetchByUserID(c, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		case errors.Is(err, domain.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

// GetAllUsers handles GET /users
// Returns all registered users
func (uc *UserController) GetAllUsers(c *gin.Context) {
	tasks, err := uc.UserUsecase.FetchAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch all users"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// Login handles POST /auth/login
// Authenticates user, sets JWT token as cookie if successful
func (uc *UserController) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, token, err := uc.UserUsecase.Login(c, body.Username, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
		case errors.Is(err, domain.ErrIncorrectPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		}
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authentication", token, 60*60*24, "", "", true, true)
	c.IndentedJSON(http.StatusOK, user)
}

// Promote handles PUT /users/:id/promote
// Promotes a user to admin by their ID
func (uc *UserController) Promote(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found"})
		return
	}

	count, err := uc.UserUsecase.PromoteByUserID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"error": "user updated successfully"})
}

// Register handles POST /auth/register
// Registers a new user and hashes their password
// The first user is automatically assigned admin rights
func (uc *UserController) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username" bson:"username" binding:"required"`
		Password string `json:"password" bson:"-" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	hashedString, err := infrastructure.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash user password"})
		return
	}

	exists, err := uc.UserUsecase.CheckIfUsernameExists(c, body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	// users, err := uc.UserUsecase.FetchAllUsers(c)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users document"})
	// 	return
	// }
	//
	// // Check for username uniqueness
	// for _, user := range users {
	// 	if user.Username == body.Username {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
	// 		return
	// 	}
	// }
	//

	count, err := uc.UserUsecase.CountUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
		return
	}

	user := domain.User{
		Username: body.Username,
		Password: hashedString,
		IsAdmin:  count == 0, // first registered user becomes admin
	}

	err = uc.UserUsecase.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert user document"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"data":    user,
	})
}
