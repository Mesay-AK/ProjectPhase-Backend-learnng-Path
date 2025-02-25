package controllers

import (
	"net/http"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

// Create handles the creation of a new task.
//
// It binds the request body to a Task struct and validates the format.
// If the format is invalid, it returns a JSON response with a "invalid task format" message and a status code of 400.
//
// It generates a new ObjectID for the task and calls the CreateTask method of the TaskUsecase to create the task.
// If there is an error during the creation process, it returns a JSON response with the error message and a status code of 500.
//
// If the task is created successfully, it returns a JSON response with the created task and a status code of 201.
func (tc *TaskController) Create(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid task format"})
		return
	}

	task.ID = primitive.NewObjectID()
	err = tc.TaskUsecase.CreateTask(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Fetch retrieves a list of tasks from the task use case and returns them as a JSON response.
// If an error occurs during the retrieval process, it returns an internal server error with the error message.
func (tc *TaskController) Fetch(c *gin.Context) {
	tasks, err := tc.TaskUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// UpdateTask updates a task based on the provided task data and task ID.
// It expects a JSON payload containing the updated task data in the request body.
// The task ID is extracted from the request URL parameters.
// If the task data is invalid or the task ID is missing, it returns a JSON response with an appropriate error message.
// If the task update is successful, it returns a JSON response with the updated task data.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	var task domain.TaskUpdate
	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid task format"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id param missing"})
		return
	}

	updatedTask, err := tc.TaskUsecase.UpdateTask(c, &task, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, updatedTask)
}

// GetTaskByID retrieves a task by its ID.
// It takes a gin.Context as a parameter and uses the "id" parameter from the context to fetch the task.
// If the "id" parameter is missing, it returns a JSON response with a status code of http.StatusBadRequest and a message indicating that the "id" parameter is missing.
// If the task with the given ID is not found, it returns a JSON response with a status code of http.StatusNotFound and a message indicating that the task with the given ID was not found.
// If there is an error while fetching the task, it returns a JSON response with a status code of http.StatusInternalServerError and the error message.
// Otherwise, it returns a JSON response with a status code of http.StatusFound and the retrieved task.
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id param missing"})
		return
	}

	task, err := tc.TaskUsecase.GetByID(c, taskID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "task with the given id not found"})
		return
	}

	c.JSON(http.StatusFound, task)
}

// DeleteTask deletes a task based on the provided task ID.
// It takes a gin.Context object as a parameter and returns no values.
// If the task ID is missing, it will respond with a JSON message indicating the missing parameter.
// If an error occurs during the deletion process, it will respond with a JSON error message.
// If the task is deleted successfully, it will respond with a JSON message indicating the successful deletion.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id param missing"})
		return
	}

	err := tc.TaskUsecase.DeleteTask(c, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task removed successfully."})
}
