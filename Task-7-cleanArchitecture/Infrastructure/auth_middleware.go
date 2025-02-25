package Infrastructure

import (
	"net/http"
	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"
	"fmt"
	"strings"
	"os"
)


func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		clientToken := ctx.GetHeader("Authorization")
		if clientToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			ctx.Abort()
			return
		}

	
		claims, err := ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			ctx.Abort()
			return
		}

		// Set claims in the context
		ctx.Set("email", claims.Email)
		ctx.Set("name", claims.Name)
		ctx.Set("uid", claims.Uid)
		ctx.Set("user_type", claims.User_type)
		ctx.Next()
	}
}

func AuthorizeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		
		authHeader := c.GetHeader("Authorization")
		var jwtSecret = []byte(os.Getenv("SECRET_KEY"))
		authParts := strings.Split(authHeader, " ")
		
		token, _:= jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		
			return jwtSecret, nil
		})
		
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		role, ok := claims["Role"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid JWT role"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(403, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
