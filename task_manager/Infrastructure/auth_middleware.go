package infrastructure

import (
	"context"
	"errors"
	"net/http"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"github.com/gin-gonic/gin"
)

// AuthenticatedUser represents a user extracted from a validated JWT
type AuthenticatedUser struct {
	ID       string
	Username string
	IsAdmin  bool
}

// AuthenticationMiddleware validates the JWT from the "Authentication" cookie,
// fetches the user from the database, and attaches it to the request context
func AuthenticationMiddleware(userRepo domain.UserRepository, jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token
		tokenString, err := c.Cookie("Authentication")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			return
		}

		// Validate JWT token
		claims, err := jwtService.Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Check claims
		if claims.Username == "" || claims.Subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing claim ID or Username"})
			return
		}

		// Check user exists in DB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userRepo.FetchByUserID(ctx, claims.Subject)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidTaskID):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			case errors.Is(err, domain.ErrUserNotFound):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user does not exist"})
			case errors.Is(err, context.DeadlineExceeded):
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "request context expired"})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			}
			return
		}

		// Set authenticated user into context
		c.Set("user", AuthenticatedUser{
			ID:       user.ID,
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		})
		c.Next()
	}
}

// AuthorizationMiddleware ensures that the authenticated user is an admin
func AuthorizationMiddleware(userRepo domain.UserRepository, jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from context (must be set by AuthenticationMiddleware)
		u, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "missing user in context"})
			return
		}

		// Assert the user context type
		userCtx, ok := u.(AuthenticatedUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
			return
		}

		// Check the user's latest admin status from DB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userRepo.FetchByUserID(ctx, userCtx.ID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrUserNotFound):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user does not exist"})
			case errors.Is(err, context.DeadlineExceeded):
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "request context expired"})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			}
			return
		}

		// Check admin privileges
		if !user.IsAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user must be admin"})
			return
		}

		// Update context user as admin and proceed
		userCtx.IsAdmin = true
		c.Set("user", userCtx)
		c.Next()
	}
}
