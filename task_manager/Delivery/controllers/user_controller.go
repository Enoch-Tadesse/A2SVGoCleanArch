package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	domain "github.com/A2SVTask7/Domain"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	repositories "github.com/A2SVTask7/Repositories"
	"github.com/gin-gonic/gin"
)

// UserController handles user-related HTTP endpoints
type UserController struct {
	UserUsecase domain.UserUsercase
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
		if errors.Is(err, repositories.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
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
	if len(tasks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there is no user yet"})
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

	user, err := uc.UserUsecase.FetchByUsername(c, body.Username)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user id"})
		return
	}

	err = infrastructure.ComparePassword(user.Password, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		return
	}

	// Build JWT claims
	claims := map[string]any{
		"sub":      user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	js := infrastructure.NewJWTService(os.Getenv("JWT_SECRET"))
	token, err := js.Generate(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to generate jwt token"})
		return
	}

	// Set JWT token as cookie
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

	users, err := uc.UserUsecase.FetchAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users document"})
		return
	}

	// Check for username uniqueness
	for _, user := range users {
		if user.Username == body.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
	}

	user := domain.User{
		Username: body.Username,
		Password: hashedString,
		IsAdmin:  len(users) == 0, // first registered user becomes admin
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
