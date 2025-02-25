package middleware

import(
	"net/http"
	"task_manager_with_jwt/data"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc{
	return func(ctx *gin.Context){

		clientToken := ctx.Request.Header.Get("token")

		if clientToken == ""{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization Denied" })
			ctx.Abort()
			return
		}

		claims, err := data.ValidateToken(clientToken) 
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err })
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("name", claims.Name)
		ctx.Set("uid", claims.Uid)
		ctx.Set("user_type", claims.User_type)
		ctx.Next()
	}
}

func AuthorizeAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        userType, ok := c.Get("user_type") // Assuming `user_type` is stored in the context
        if !ok || userType != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admins only"})
            c.Abort()
            return
        }
        c.Next()
    }
}