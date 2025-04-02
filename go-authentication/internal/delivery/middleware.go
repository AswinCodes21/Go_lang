package delivery

import (
	"net/http"
	"strings"

	"go-authentication/pkg"

	"log"

	"github.com/gin-gonic/gin"
)

// ErrorResponse function defines the standard error response structure

type ErrorResponse struct {
	Message string `json:"message"`
}

//ErrorHandler Fucntion handles error globally

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
			c.Abort()
		}

	}
}

//AuthMiddleware Function validates the JWT token and extracts user information

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			c.Abort()
			return
		}

		// Debug logging
		log.Printf("Received Authorization header: %s", authHeader)

		// Extract token (format: "Bearer <token>")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Debug logging
		log.Printf("Extracted token: %s", tokenString)

		// Validate token and get claims
		claims, err := pkg.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Debug logging
		log.Printf("Token validated successfully. User ID: %v", claims["user_id"])

		// Set user ID and email in context
		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Next()
	}
}
