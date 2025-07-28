package routers

import (
	"time"

	"github.com/A2SVTask7/Delivery/controllers"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	repositories "github.com/A2SVTask7/Repositories"
	usecases "github.com/A2SVTask7/Usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// newTaskRouter sets up routes for task operations accessible by authenticated users
func newTaskRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup, config infrastructure.Config) {
	tr := repositories.NewTaskRepository(db, config.CollectionTask)
	tc := &controllers.TaskController{
		TaskUsecase: usecases.NewTaskUsecase(tr, timeout),
	}
	group.GET("/tasks", tc.GetAllTasks)
	group.GET("/tasks/:id", tc.GetTaskByID)
}

// newUserRouter sets up public routes related to user authentication and registration
func newUserRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup, config infrastructure.Config) {
	ur := repositories.NewUserRepository(db, config.CollectionUser)
	jwt := infrastructure.NewJWTService(config.JWTSecret)
	pws := infrastructure.NewPasswordService()
	uc := &controllers.UserController{
		UserUsecase: usecases.NewUserUsecase(ur, jwt, pws, timeout),
	}
	group.POST("/login", uc.Login)
	group.POST("/register", uc.Register)
}

// newAdminRouter sets up routes for admin-level operations including user management and task CRUD
func newAdminRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup, config infrastructure.Config) {
	ur := repositories.NewUserRepository(db, config.CollectionUser)
	jwt := infrastructure.NewJWTService(config.JWTSecret)
	pws := infrastructure.NewPasswordService()
	uc := &controllers.UserController{
		UserUsecase: usecases.NewUserUsecase(ur, jwt, pws, timeout),
	}

	tr := repositories.NewTaskRepository(db, config.CollectionTask)
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
func SetUp(timeout time.Duration, db mongo.Database, router *gin.Engine, config infrastructure.Config) {
	userRepo := repositories.NewUserRepository(db, config.CollectionUser)
	jwtService := infrastructure.NewJWTService(config.JWTSecret) // Use environment variable for JWT secret

	authMiddleware := infrastructure.AuthenticationMiddleware(userRepo, jwtService)
	adminMiddleware := infrastructure.AuthorizationMiddleware(userRepo, jwtService)

	// Public routes without authentication
	publicRouter := router.Group("")
	newUserRouter(timeout, db, publicRouter, config)

	// Routes requiring authentication
	authenticatedRouter := router.Group("")
	authenticatedRouter.Use(authMiddleware)
	newTaskRouter(timeout, db, authenticatedRouter, config)

	// Admin-only routes require both authentication and authorization
	adminRouter := router.Group("")
	adminRouter.Use(authMiddleware, adminMiddleware)
	newAdminRouter(timeout, db, adminRouter, config)
}
