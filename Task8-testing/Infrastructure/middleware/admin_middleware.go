package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware is a middleware function that checks if the user is an admin.
// It checks for the presence of "claims" in the context and verifies that the user's role is "ADMIN".
// If the user is not an admin, it returns a JSON response with an error message and aborts the request.
// Otherwise, it allows the request to proceed to the next handler.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not allowed"})
			c.Abort()
			return
		}

		userRole, ok := c.Get("userRole")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not allowed"})
			c.Abort()
			return
		}
		roleStr, _ := userRole.(string)

		if strings.ToUpper(roleStr) != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
