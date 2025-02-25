	package router

	import(
		"task_manager/controllers"
		"github.com/gin-gonic/gin"
	)
	
	
	func SetRouter() *gin.Engine{

	router := gin.Default()

	router.GET("/tasks", controller.GetTasks)
	router.GET("/tasks/:id", controller.GetTaskById)
	router.POST("/tasks", controller.AddTask)
	router.PUT("tasks/:id",controller.UpdateTask)
	router.DELETE("tasks/:id", controller.DeleteTask)

	return router
}