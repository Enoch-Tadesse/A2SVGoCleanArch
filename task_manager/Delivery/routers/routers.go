package routers

import (
	"os"
	"time"

	"github.com/A2SVTask7/Delivery/controllers"
	domain "github.com/A2SVTask7/Domain"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	repositories "github.com/A2SVTask7/Repositories"
	usecases "github.com/A2SVTask7/Usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// newTaskRouter sets up routes for task operations accessible by authenticated users
func newTaskRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repositories.NewTaskRepository(db, domain.CollectionTask)
	tc := &controllers.TaskController{
		TaskUsecase: usecases.NewTaskUsecase(tr, timeout),
	}
	group.GET("/tasks", tc.GetAllTasks)
	group.GET("/tasks/:id", tc.GetTaskByID)
}

// newUserRouter sets up public routes related to user authentication and registration
func newUserRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)
	uc := &controllers.UserController{
		UserUsecase: usecases.NewUserUsecase(ur, timeout),
	}
	group.POST("/login", uc.Login)
	group.POST("/register", uc.Register)
}

// newAdminRouter sets up routes for admin-level operations including user management and task CRUD
func newAdminRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, domain.CollectionUser)
	uc := &controllers.UserController{
		UserUsecase: usecases.NewUserUsecase(ur, timeout),
	}

	tr := repositories.NewTaskRepository(db, domain.CollectionTask)
	tc := &controllers.TaskController{
		TaskUsecase: usecases.NewTaskUsecase(tr, timeout),
	}

	group.GET("/users", uc.GetAllUsers)
	group.GET("/users/:id", uc.GetUserByID)
	group.PATCH("/promote/:id", uc.Promote)
	group.POST("/tasks", tc.CreateTask)
	group.DELETE("/tasks/:id", tc.DeleteTask)
	group.PUT("/tasks/:id", tc.UpdateTask)
}

// SetUp configures all the route groups and applies middleware for authentication and authorization
func SetUp(timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	userRepo := repositories.NewUserRepository(db, domain.CollectionUser)
	jwtService := infrastructure.NewJWTService(os.Getenv("JWT_SECRET")) // Use environment variable for JWT secret

	// Public routes without authentication
	publicRouter := gin.Group("")
	newUserRouter(timeout, db, publicRouter)

	// Routes requiring authentication
	authenticatedRouter := gin.Group("")
	authenticatedRouter.Use(infrastructure.AuthenticationMiddleware(userRepo, jwtService))
	newTaskRouter(timeout, db, authenticatedRouter)

	// Admin-only routes require both authentication and authorization
	adminRouter := gin.Group("")
	adminRouter.Use(infrastructure.AuthenticationMiddleware(userRepo, jwtService))
	adminRouter.Use(infrastructure.AuthorizationMiddleware(userRepo, jwtService))
	newAdminRouter(timeout, db, adminRouter)
}
