package controllers

import (
	"fmt"
	"net/http"
	domain "task_manager_clean/Domain"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	TaskUsecase domain.TaskUsecase
	UserUsecase domain.UserUsecase
}


func respondWithError(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"error": message})
}

func (ctrl *Controller) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user domain.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			respondWithError(ctx, http.StatusBadRequest, "Invalid input")
			return
		}

		if err := ctrl.UserUsecase.Register(ctx, &user); err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to register user")
			fmt.Print("error out of the register of usecase, in controller")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}


func (ctrl *Controller) LogIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user domain.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			respondWithError(ctx, http.StatusBadRequest, "Invalid input")
			return
		}

		foundUser, token, refreshToken, err := ctrl.UserUsecase.LogIn(ctx, *user.Email, *user.Password)
		if err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to log in")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user":          foundUser,
			"access_token":  token,
			"refresh_token": refreshToken,
		})
	}
}

func (ctrl *Controller) GetUsers(ctx *gin.Context) {
	users, err := ctrl.UserUsecase.GetUsers(ctx)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	if len(users) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"users": users})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "No users found"})
}


func (ctrl *Controller) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := ctrl.UserUsecase.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			respondWithError(ctx, http.StatusNotFound, "User not found")
		} else {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve user")
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}


func (ctrl *Controller) PromoteUser(ctx *gin.Context) {
    id := ctx.Param("id")

	var promotedUser domain.User

    user, err := ctrl.UserUsecase.PromoteUser(ctx, id, &promotedUser)
    if err != nil {
        if err == domain.ErrUserNotFound {
            respondWithError(ctx, http.StatusNotFound, "User not found")
        } else {
            respondWithError(ctx, http.StatusInternalServerError, "Failed to promote user")
        }
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "User promoted successfully", "user": user})
}


func (ctrl *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ctrl.UserUsecase.DeleteUser(ctx, id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			respondWithError(ctx, http.StatusNotFound, "User not found")
		} else {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctrl *Controller) AddTask(ctx *gin.Context) {
	var task domain.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := ctrl.TaskUsecase.AddTask(ctx, &task); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to create task")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}


func (ctrl *Controller) GetTasks(ctx *gin.Context) {
	tasks, err := ctrl.TaskUsecase.GetTasks(ctx)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	if len(tasks) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
}

func (ctrl *Controller) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := ctrl.TaskUsecase.GetTaskByID(ctx, id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			respondWithError(ctx, http.StatusNotFound, "Task not found")
		} else {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve task")
		}
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (ctrl *Controller) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask domain.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	task, err := ctrl.TaskUsecase.UpdateTask(ctx, id, &updatedTask)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			respondWithError(ctx, http.StatusNotFound, "Task not found")
		} else {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to update task")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}


func (ctrl *Controller) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ctrl.TaskUsecase.DeleteTask(ctx, id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			respondWithError(ctx, http.StatusNotFound, "Task not found")
		} else {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to delete task")
		}

	}

}