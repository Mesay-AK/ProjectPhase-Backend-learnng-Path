package controllers

import (
	"net/http"
	"task_manager_with_jwt/data"
	"task_manager_with_jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Controller struct {
	TaskService *data.TaskService
	UserService *data.UserService
}


func (ctrl *Controller) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Register the user
		if err := ctrl.UserService.Register(ctx, &user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User registerd successfully"})
	}
}

func (ctrl *Controller) LogIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		foundUser, token, refreshToken, err := ctrl.UserService.LogIn(*user.Email, *user.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	users, err := ctrl.UserService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Users"})
		return
	}

	if len(users) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"Users": users})
		return
	} 
	ctx.JSON(http.StatusOK, gin.H{"message": "No user found"})
}

func (ctrl *Controller)GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}


func (ctrl *Controller) PromoteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var promotedUser models.User

	if err := ctx.ShouldBindJSON(&promotedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserService.PromoteUser(id, promotedUser)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User Promoted successfully", "task": user})
}

func (ctrl *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ctrl.UserService.DeleteUser(id)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete User"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}















/*
Accessing the task collection and manipulating it.
*/


func (ctrl *Controller) AddTask(ctx *gin.Context) {
	var task models.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.TaskService.AddTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}
func (ctrl *Controller) GetTasks(ctx *gin.Context) {
	tasks, err := ctrl.TaskService.GetTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	if len(tasks) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
		return
	} 
	ctx.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
}

func (ctrl *Controller)GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := ctrl.TaskService.GetTaskByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (ctrl *Controller) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := ctrl.TaskService.UpdateTask(id, updatedTask)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}

func (ctrl *Controller) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ctrl.TaskService.DeleteTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}