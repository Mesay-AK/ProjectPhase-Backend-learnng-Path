package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskControllerTest struct {
	suite.Suite
	usecase       *mocks.TaskUsecase
	controller    *TaskController
	testingServer *httptest.Server
	timeout       time.Duration
	nilTask       domain.Task
}

var updatedTitle = "Updated Task"
var updateDescription = "Updated description"
var updateStatus = domain.TaskStatus("In Progress")
var updateDate = time.Now()

func (suite *taskControllerTest) SetupSuite() {
	suite.usecase = new(mocks.TaskUsecase)
	suite.controller = &TaskController{
		TaskUsecase: suite.usecase,
	}
	suite.timeout = time.Duration(10) * time.Second
	suite.nilTask = domain.Task{}
	router := gin.Default()
	router.POST("/tasks", suite.controller.Create)
	router.GET("/tasks", suite.controller.Fetch)
	router.PUT("/tasks/:id", suite.controller.UpdateTask)
	router.GET("/tasks/:id", suite.controller.GetTaskByID)

	testingServer := httptest.NewServer(router)
	suite.testingServer = testingServer
}

func (suite *taskControllerTest) TearDownSuite() {
	defer suite.testingServer.Close()
}

func (suite *taskControllerTest) TestCreateTask() {
	// Mock request payload
	requestPayload := domain.Task{
		Title:       "New Task",
		Description: "This is a task description",
		Status:      domain.TaskStatus("Started"),
		Due_Date:    time.Now(),
	}

	// Convert request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	suite.NoError(err)

	// Mock the usecase response
	suite.usecase.On("CreateTask", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.Task)
		arg.ID = primitive.NewObjectID()
	}).Once()

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/tasks", bytes.NewBuffer(requestBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusCreated, resp.StatusCode)

	// Decode the response body
	var responseBody domain.Task
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	// Assert the response body
	suite.NotEqual(primitive.NilObjectID, responseBody.ID)
	suite.Equal(requestPayload.Title, responseBody.Title)
	suite.Equal(requestPayload.Description, responseBody.Description)
	suite.Equal(requestPayload.Status, responseBody.Status)
	suite.WithinDuration(requestPayload.Due_Date, responseBody.Due_Date, time.Second)

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerTest) TestCreateInvalidTask() {
	// Mock invalid request payload (missing Title)
	requestPayload := domain.Task{
		Description: "This is a task description",
		Status:      domain.TaskStatus("Started"),
		Due_Date:    time.Now(),
	}

	// Convert request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	suite.NoError(err)

	// Mock the usecase response to return an error
	suite.usecase.On("CreateTask", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(errors.New("invalid task")).Once()

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/tasks", bytes.NewBuffer(requestBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())

	// Decode the response body to check the error message
	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)

	// Assert the error message
	suite.Equal("invalid task", responseBody["error"])
}

func (suite *taskControllerTest) TestFetch() {
	// Mock the usecase response to return a list of tasks
	expectedTasks := []domain.Task{
		{ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: domain.TaskStatus("Started"), Due_Date: time.Now()},
		{ID: primitive.NewObjectID(), Title: "Task 2", Description: "Description 2", Status: domain.TaskStatus("Completed"), Due_Date: time.Now()},
	}
	suite.usecase.On("Fetch", mock.Anything).Return(expectedTasks, nil).Once()

	// Create a new HTTP request to fetch tasks
	req, err := http.NewRequest(http.MethodGet, suite.testingServer.URL+"/tasks", nil)
	suite.NoError(err)

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusOK, resp.StatusCode)

	// Decode the response body and assert the tasks
	var responseBody []domain.Task
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)
	suite.Equal(len(expectedTasks), len(responseBody))
	for i, task := range responseBody {
		suite.Equal(expectedTasks[i].ID, task.ID)
		suite.Equal(expectedTasks[i].Title, task.Title)
		suite.Equal(expectedTasks[i].Description, task.Description)
		suite.Equal(expectedTasks[i].Status, task.Status)
		suite.WithinDuration(expectedTasks[i].Due_Date, task.Due_Date, time.Second)
	}

	// Mock the usecase response to return an error
	suite.usecase.On("Fetch", mock.Anything).Return(nil, errors.New("fetch error")).Once()

	// Create a new HTTP request to fetch tasks
	req, err = http.NewRequest(http.MethodGet, suite.testingServer.URL+"/tasks", nil)
	suite.NoError(err)

	// Send the request and get the response
	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)

	// Decode the response body and assert the error message
	var errorResponseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResponseBody)
	suite.NoError(err)
	suite.Equal("fetch error", errorResponseBody["error"])

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())
}
func (suite *taskControllerTest) TestUpdateTask() {
	// Mock request payload and task ID
	updatedTitle := "Updated Task"
	updateDescription := "Updated description"
	updateStatus := domain.TaskStatus("In Progress")
	updateDate := time.Now()
	requestPayload := domain.TaskUpdate{
		Title:       &updatedTitle,
		Description: &updateDescription,
		Status:      &updateStatus,
		Due_Date:    &updateDate,
	}
	taskID := "some-task-id"

	// Convert request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	suite.NoError(err)

	// Mock the usecase response to return the updated task
	expectedTask := domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       *requestPayload.Title,
		Description: *requestPayload.Description,
		Status:      *requestPayload.Status,
		Due_Date:    *requestPayload.Due_Date,
	}
	suite.usecase.On("UpdateTask", mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(taskUpdate *domain.TaskUpdate) bool {
		return *taskUpdate.Title == *requestPayload.Title &&
			*taskUpdate.Description == *requestPayload.Description &&
			*taskUpdate.Status == *requestPayload.Status &&
			taskUpdate.Due_Date.Equal(*requestPayload.Due_Date)
	}), taskID).Return(expectedTask, nil).Once()

	// Create a new HTTP request to update the task
	req, err := http.NewRequest(http.MethodPut, suite.testingServer.URL+"/tasks/"+taskID, bytes.NewBuffer(requestBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusCreated, resp.StatusCode)

	// Decode the response body and assert the updated task
	var responseBody domain.Task
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)
	suite.Equal(expectedTask.ID, responseBody.ID)
	suite.Equal(expectedTask.Title, responseBody.Title)
	suite.Equal(expectedTask.Description, responseBody.Description)
	suite.Equal(expectedTask.Status, responseBody.Status)
	suite.WithinDuration(expectedTask.Due_Date, responseBody.Due_Date, time.Second)

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerTest) TestUpdateTaskInvalidFormat() {
	// Mock invalid request payload
	requestPayload := map[string]interface{}{
		"Title": 123, // Invalid type
	}

	// Convert request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	suite.NoError(err)

	// Create a new HTTP request to update the task
	req, err := http.NewRequest(http.MethodPut, suite.testingServer.URL+"/tasks/some-task-id", bytes.NewBuffer(requestBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	// Decode the response body and assert the error message
	var errorResponseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResponseBody)
	suite.NoError(err)
	suite.Equal("invalid task format", errorResponseBody["message"])

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerTest) TestUpdateTaskMissingID() {
	// Mock request payload
	requestPayload := domain.TaskUpdate{
		Title:       &updatedTitle,
		Description: &updateDescription,
		Status:      &updateStatus,
		Due_Date:    &updateDate,
	}

	// Convert request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	suite.NoError(err)

	// Create a new HTTP request to update the task without task ID
	req, err := http.NewRequest(http.MethodPut, suite.testingServer.URL+"/tasks/", bytes.NewBuffer(requestBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusNotFound, resp.StatusCode) // Adjusted to 404

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerTest) TestDeleteTaskSuccess() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	taskID := "some-task-id"
	req, err := http.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
	suite.NoError(err)
	c.Params = gin.Params{gin.Param{Key: "id", Value: taskID}}
	c.Request = req

	suite.usecase.On("DeleteTask", mock.Anything, taskID).Return(nil).Once()

	tc := &TaskController{TaskUsecase: suite.usecase}
	tc.DeleteTask(c)

	suite.Equal(http.StatusOK, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("task removed successfully.", response["message"])

	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerTest) TestGetTaskByIDSuccess() {
	taskID := "some-task-id"

	// Mock the usecase response to return the expected task
	expectedTask := domain.Task{
		ID:    primitive.NewObjectID(),
		Title: "Test Task",
	}
	suite.usecase.On("GetByID", mock.AnythingOfType("*gin.Context"), taskID).Return(expectedTask, nil).Once()

	// Create a new HTTP request with task ID
	req, err := http.NewRequest(http.MethodGet, suite.testingServer.URL+"/tasks/"+taskID, nil)
	suite.NoError(err)

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusFound, resp.StatusCode)

	// Decode the response body and assert the retrieved task
	var responseBody domain.Task
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	suite.NoError(err)
	suite.Equal(expectedTask.ID, responseBody.ID)
	suite.Equal(expectedTask.Title, responseBody.Title)

	// Assert the mock expectations
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerTest) TestGetTaskByIDMissingID() {
	// Define the expected error message
	expectedErrorMessage := "id param missing"

	// Setup the mock to return an error when the ID is missing
	suite.usecase.On("Fetch", mock.Anything).Return(nil, errors.New(expectedErrorMessage))

	// Create a new HTTP request without a task ID
	req, err := http.NewRequest(http.MethodGet, suite.testingServer.URL+"/tasks/", nil)
	suite.NoError(err)

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Assert the response status code
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(taskControllerTest))
}
