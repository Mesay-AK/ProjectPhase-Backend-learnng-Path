	package router

	import(
		"task_manager_using_mongoDB/controllers"
		"github.com/gin-gonic/gin"
	)
	
	
	func SetRouter(taskController *controllers.TaskController) *gin.Engine{

	router := gin.Default()

	router.GET("/tasks", taskController.GetTasks)
	router.GET("/tasks/:id", taskController.GetTaskById)
	router.POST("/tasks", taskController.AddTask)
	router.PUT("tasks/:id",taskController.UpdateTask)
	router.DELETE("tasks/:id", taskController.DeleteTask)

	return router
}