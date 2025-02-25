package main

import (
	"time"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Delivery/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.Close()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	routers.Setup(env, timeout, *db, gin)

	gin.Run(env.ServerAddress)
}
