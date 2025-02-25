package controller

import (
	// "fmt"
	"net/http"
	"task_manager/models"
	"task_manager/data"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context){
	tasks:= data.GetTasks()

	if len(tasks) > 0 {
		ctx.JSON(200, gin.H{"tasks": tasks,})
		return
	} 
	ctx.JSON(200, gin.H{"tasks": "Zero tasks found",})

	}
func GetTaskById(ctx *gin.Context){

		id := ctx.Param("id")

		task, err := data.GetTaskByID(id);
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error:":"Task not found"})
			return
		}
		ctx.JSON(http.StatusOK, task)
		
	}

func AddTask(ctx *gin.Context) {
	var task models.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.AddTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func UpdateTask(ctx *gin.Context){
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err!= nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.UpdateTask(id, updatedTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task Updated"})
}

func DeleteTask(ctx *gin.Context){
	id := ctx.Param("id")

	err := data.DeleteTask(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task Deleted"})

}


