package routers

import (
	"time"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Delivery/controllers"
	repositories "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Repositories"
	usecases "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewTaskAdminRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repositories.NewTaskRepository(db, "tasks")
	tc := controllers.TaskController{
		TaskUsecase: usecases.NewTaskUsecase(tr, timeout),
	}
	group.POST("/tasks", tc.Create)
	group.PUT("/tasks/:id", tc.UpdateTask)
	group.DELETE("/tasks/:id", tc.DeleteTask)
}
