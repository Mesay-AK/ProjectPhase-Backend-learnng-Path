package routers

import (
	"task_manager_clean/Delivery/controllers"
	infrastructure "task_manager_clean/Infrastructure"
	"github.com/gin-gonic/gin"
)
func SetRouter(controller *controllers.Controller) *gin.Engine {
    router := gin.Default()

   
    router.POST("/register", controller.SignUp()) 
    router.POST("/login", controller.LogIn()) 
  

    userRoutes := router.Group("/user")
    userRoutes.Use(infrastructure.Authenticate()) 
    {
        userRoutes.GET("/tasks", controller.GetTasks) 
        userRoutes.GET("/tasks/:id", controller.GetTaskByID)
        
    }

    adminRoutes := router.Group("/admin")
    adminRoutes.Use(infrastructure.Authenticate(), infrastructure.AuthorizeAdmin()) 
    {
        adminRoutes.GET("/tasks", controller.GetTasks)
        adminRoutes.POST("/tasks", controller.AddTask)
        adminRoutes.DELETE("/tasks/:id", controller.DeleteTask)
        adminRoutes.PUT("/tasks/:id", controller.UpdateTask)
        adminRoutes.PATCH("/promote/:id", controller.PromoteUser) 
        userRoutes.GET("/users", controller.GetUsers)
        userRoutes.GET("/users/:id", controller.GetUserByID)
        adminRoutes.DELETE("/users/:id", controller.DeleteUser)
    }

    return router
}