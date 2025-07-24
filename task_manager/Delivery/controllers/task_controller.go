package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"github.com/gin-gonic/gin"
)

// TaskController handles incoming HTTP requests related to tasks
type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

// CreateTask handles POST /tasks
// Validates the request, checks for due date, creates a new task
func (tc *TaskController) CreateTask(c *gin.Context) {
	log.Println("herer")
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

	if err := tc.TaskUsecase.Create(c, &task); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDueDate):
			c.JSON(http.StatusBadRequest, gin.H{"error": "due date can't be in the past"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		}
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
	err := tc.TaskUsecase.DeleteByTaskID(c, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		}
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
		if errors.Is(err, domain.ErrTaskNotFound) {
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

	task := domain.Task{
		ID:          id,
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
		Status:      body.Status,
	}

	err := tc.TaskUsecase.UpdateByTaskID(c, &task)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		case errors.Is(err, domain.ErrInvalidDueDate):
			c.JSON(http.StatusBadRequest, gin.H{"error": "due date can not be in the past"})
		case errors.Is(err, domain.ErrTaskNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "task not found"})
		case errors.Is(err, domain.ErrNoChangesMade):
			c.JSON(http.StatusOK, gin.H{
				"message": "no changes were made",
				"data":    task,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}
