package main

import(
	"task_manager/router"
	"log"
)

func main(){

	Router := router.SetRouter()

	if err := Router.Run(":8080"); err != nil {

		log.Fatalf("Failed to start server: %v", err)
	}
}