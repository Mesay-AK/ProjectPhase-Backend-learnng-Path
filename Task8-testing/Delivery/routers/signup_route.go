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

func NewSignupRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, "users")
	sc := controllers.SignupController{
		SignupUsecase: usecases.NewSignupUsecase(ur, timeout),
		Env:           env,
	}
	group.POST("/user/register", sc.Signup)
}
