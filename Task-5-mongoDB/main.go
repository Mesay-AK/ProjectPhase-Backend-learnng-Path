package main

import(
	"log"
	"task_manager_using_mongoDB/database"
	"task_manager_using_mongoDB/controllers"
	"task_manager_using_mongoDB/data"
	"task_manager_using_mongoDB/router"
)

func main(){

	// Setting up the database
	db, err := database.InitDataBase("mongodb://localhost:27017", "TaskDB")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	taskService := data.NewTaskService(db)
	taskController := controllers.NewTaskController(taskService)

	// Setting up the router
	Router := router.SetRouter(taskController)
	if err := Router.Run(":8080"); err != nil {

		log.Fatalf("Failed to start server: %v", err)
	}
}


