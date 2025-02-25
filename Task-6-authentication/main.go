package main

import (
	"task_manager_with_jwt/controllers"
	"task_manager_with_jwt/database"
	"task_manager_with_jwt/data"
	"task_manager_with_jwt/router"
	"log"
	"os"
	

)

func main(){
	
	// Initializing the database
	db, err := database.InitDataBase("mongodb://localhost:27017", "TaskDB")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}


	
	taskService := data.NewTaskService(db)
	userService := data.NewUserService(db)


	ctrls := &controllers.Controller{
		TaskService: taskService, 
		UserService: userService,
	}

	// Setting and running the server
	Router := router.SetRouter(ctrls)

	port := os.Getenv("PORT")

	if port == ""{
		port = ":8080"
	}
	if err := Router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}


}