package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	domain "github.com/A2SVTask7/Domain"
	repositories "github.com/A2SVTask7/Repositories"
	"github.com/gin-gonic/gin"
)

// TaskController handles incoming HTTP requests related to tasks
type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

// CreateTask handles POST /tasks
// Validates the request, checks for due date, creates a new task
func (tc *TaskController) CreateTask(c *gin.Context) {
	var body struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date" binding:"required"`
		Status      string    `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}

	task := domain.Task{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
		Status:      body.Status,
	}

	if body.DueDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due date can not be in the past"})
		return
	}

	err := tc.TaskUsecase.Create(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create task: %s", err)})
		return
	}

	c.IndentedJSON(http.StatusCreated, task)
}

// DeleteTask handles DELETE /tasks/:id
// Validates the ID and deletes the corresponding task
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id can not be empty"})
		return
	}

	// Delete the task
	count, err := tc.TaskUsecase.DeleteByTaskID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task delete successfully"})
}

// GetAllTasks handles GET /tasks
// Fetches and returns all tasks
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.TaskUsecase.FetchAllTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch all tasks"})
		return
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there is no task yet"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// GetTaskByID handles GET /tasks/:id
// Validates the ID and fetches the specific task
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	task, err := tc.TaskUsecase.FetchByTaskID(c, id)
	if err != nil {
		if errors.Is(err, repositories.ErrTaskNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task"})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// UpdateTask handles PUT /tasks/:id
// Validates the ID and request body, updates the task fields
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id can not be empty"})
		return
	}
	var body struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date" binding:"required"`
		Status      string    `json:"status" binding:"required,oneof=pending completed missed"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Normalize and validate status
	body.Status = strings.TrimSpace(body.Status)
	body.Status = strings.ToLower(body.Status)

	// Validate due date
	if body.DueDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due date can not be in the past"})
		return
	}

	task := domain.Task{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
		Status:      body.Status,
	}

	matched, modified, err := tc.TaskUsecase.UpdateByTaskID(c, id, &task)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		}
		return
	}

	if matched == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not found"})
		return
	}

	if modified == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "no changes were made",
			"data":    task,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}
