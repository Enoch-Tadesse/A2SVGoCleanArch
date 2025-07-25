package tasks

import (
	"testing"

	"github.com/A2SVTask7/Delivery/controllers"
	mock "github.com/A2SVTask7/tests/controllers/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SuiteTaskUsecase struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *mock.MockTaskUsecase
}

func (s *SuiteTaskUsecase) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUsecase = new(mock.MockTaskUsecase)
	s.router = gin.Default()
	s.router.RedirectTrailingSlash = false

	taskController := controllers.TaskController{TaskUsecase: s.mockUsecase}
	s.router.POST("/tasks", taskController.CreateTask)
	s.router.GET("/tasks", taskController.GetAllTasks)
	s.router.GET("/tasks/:id", taskController.GetTaskByID)
	s.router.PUT("/tasks/:id", taskController.UpdateTask)
}

func (s *SuiteTaskUsecase) PrepareTest(tt TaskListTestCase) {
	// Clear previous mock calls and expectations
	s.mockUsecase.ExpectedCalls = nil
	s.mockUsecase.Calls = nil

	// If there's a MockSetup function, run it
	if tt.MockSetup != nil {
		tt.MockSetup()
	}
}

func TestTaskController(t *testing.T) {
	suite.Run(t, new(SuiteTaskUsecase))
}
