package users

import (
	"testing"

	"github.com/A2SVTask7/Delivery/controllers"
	mock "github.com/A2SVTask7/tests/controllers_test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SuiteUserUsecase struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *mock.MockUserUsecase
}

func (s *SuiteUserUsecase) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUsecase = new(mock.MockUserUsecase)
	s.router = gin.Default()
	s.router.RedirectTrailingSlash = false

	taskController := controllers.UserController{UserUsecase: s.mockUsecase}
	s.router.GET("/users", taskController.GetAllUsers)
	s.router.POST("/users", taskController.Register)
	s.router.POST("/login", taskController.Login)
	s.router.GET("/users/:id", taskController.GetUserByID)
	s.router.PATCH("/promote/:id", taskController.Promote)
}

func (s *SuiteUserUsecase) PrepareTest(tt UserListTestCase) {
	// Clear previous mock calls and expectations
	s.mockUsecase.ExpectedCalls = nil
	s.mockUsecase.Calls = nil

	// If there's a MockSetup function, run it
	if tt.MockSetup != nil {
		tt.MockSetup()
	}
}

func TestTaskController(t *testing.T) {
	suite.Run(t, new(SuiteUserUsecase))
}
