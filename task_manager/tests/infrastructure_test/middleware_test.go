package infrastructure_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/A2SVTask7/Domain"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	"github.com/A2SVTask7/tests/infrastructure_test/mocks"
	mockRepo "github.com/A2SVTask7/tests/usecases_test"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockJWT  *mocks.MockJWTService
	mockRepo *mockRepo.MockUserRepository
}

var fakeClaims = &infrastructure.CustomClaims{
	Username: "testuser",
	RegisteredClaims: jwt.RegisteredClaims{
		Subject:   "user-id-123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	},
}

var sampleUser = domain.User{
	ID:       "user-id-123",
	Username: "testuser",
	IsAdmin:  true,
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockJWT = new(mocks.MockJWTService)
	suite.mockRepo = new(mockRepo.MockUserRepository)
	suite.router = gin.New()
}

func (suite *AuthMiddlewareTestSuite) TestAuthenticationMiddleware_Success() {
	suite.mockJWT.On("Validate", "valid-token").Return(fakeClaims, nil)
	suite.mockRepo.On("FetchByUserID", mock.Anything, "user-id-123").Return(sampleUser, nil)

	suite.router.Use(infrastructure.AuthenticationMiddleware(suite.mockRepo, suite.mockJWT))
	suite.router.GET("/protected", func(c *gin.Context) {
		u, exists := c.Get("user")
		suite.True(exists)
		user := u.(infrastructure.AuthenticatedUser)
		suite.Equal("testuser", user.Username)
		suite.True(user.IsAdmin)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "Authentication", Value: "valid-token"})
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.mockJWT.AssertExpectations(suite.T())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestAuthenticationMiddleware_InvalidToken() {
	suite.mockJWT.On("Validate", "bad-token").Return(nil, assert.AnError)

	suite.router.Use(infrastructure.AuthenticationMiddleware(suite.mockRepo, suite.mockJWT))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "Authentication", Value: "bad-token"})
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.mockJWT.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestAuthorizationMiddleware_Success() {
	sampleUser.IsAdmin = true
	suite.mockRepo.On("FetchByUserID", mock.Anything, "user-id-123").Return(sampleUser, nil)

	suite.router.Use(func(c *gin.Context) {
		c.Set("user", infrastructure.AuthenticatedUser{
			ID:       "user-id-123",
			Username: "testuser",
		})
		c.Next()
	})
	suite.router.Use(infrastructure.AuthorizationMiddleware(suite.mockRepo, suite.mockJWT))

	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin router"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestAuthorizationMiddleware_NotAdmin() {
	sampleUser.IsAdmin = false
	suite.mockRepo.On("FetchByUserID", mock.Anything, "user-id-123").Return(sampleUser, nil)

	suite.router.Use(func(c *gin.Context) {
		c.Set("user", infrastructure.AuthenticatedUser{
			ID:       "user-id-123",
			Username: "testuser",
		})
		c.Next()
	})
	suite.router.Use(infrastructure.AuthorizationMiddleware(suite.mockRepo, suite.mockJWT))

	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin router"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
