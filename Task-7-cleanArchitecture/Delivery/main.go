package main

import (
	"log"
	"os"
	controller "task_manager_clean/Delivery/controllers"
	router "task_manager_clean/Delivery/routers"
	infrastructure "task_manager_clean/Infrastructure"
	repository "task_manager_clean/Repositories"
	usecase "task_manager_clean/Usecases"
	"github.com/go-playground/validator/v10"

)

func main() {

	db, err := infrastructure.InitDataBase("mongodb://localhost:27017", "TaskDB")
	
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	validate := validator.New()

	taskUsecase := usecase.NewTaskUseCase(taskRepo, validate)
	userUsecase := usecase.NewUserUseCase(userRepo, validate)

	ctrls := &controller.Controller{
		TaskUsecase: taskUsecase,
		UserUsecase: userUsecase,
	}

	Router := router.SetRouter(ctrls)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	if err := Router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
