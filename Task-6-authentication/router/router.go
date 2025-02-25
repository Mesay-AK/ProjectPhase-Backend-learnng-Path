package router

import (
	"task_manager_with_jwt/controllers"
	"task_manager_with_jwt/middleware"
	"github.com/gin-gonic/gin"
)

func SetRouter(controller *controllers.Controller) *gin.Engine {
	router := gin.Default()

	// Routes for new user

	router.POST("/register", controller.SignUp())
	router.POST("/login", controller.LogIn())

	// Routes for users
	userRoutes := router.Group("/user")
	userRoutes.Use(middleware.Authenticate())
	{
		userRoutes.GET("/tasks", controller.GetTasks) 
		userRoutes.GET("/tasks/:id", controller.GetTaskByID)  
		userRoutes.GET("/users/:id", controller.GetUserByID)
	}

	// Routes for admins
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.Authenticate(), middleware.AuthorizeAdmin())
	{
		adminRoutes.GET("/tasks", controller.GetTasks) 
		adminRoutes.POST("/tasks", controller.AddTask)
		adminRoutes.DELETE("/tasks/:id", controller.DeleteTask) 
		adminRoutes.PUT("/tasks/:id", controller.UpdateTask) 
		adminRoutes.PATCH("/promote/:id", controller.PromoteUser) 
		adminRoutes.DELETE("/users/:id", controller.DeleteUser)  
	}

	return router
}
