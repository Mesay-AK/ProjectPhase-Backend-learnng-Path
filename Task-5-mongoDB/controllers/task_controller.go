package controllers

import (
	"net/http"
	"task_manager_using_mongoDB/models"
	"task_manager_using_mongoDB/data"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskService *data.TaskService
}

func NewTaskController(taskService *data.TaskService) *TaskController {
	return &TaskController{TaskService: taskService}
}

func (c *TaskController) AddTask(ctx *gin.Context) {
	var task models.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.TaskService.AddTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func (c *TaskController) GetTasks(ctx *gin.Context) {
	tasks, err := c.TaskService.GetTasks()
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

func (c *TaskController)GetTaskById(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := c.TaskService.GetTaskByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := c.TaskService.UpdateTask(id, updatedTask)
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

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.TaskService.DeleteTask(id)
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
