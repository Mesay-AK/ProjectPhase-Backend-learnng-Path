package usecases

import (
	"context"
	"testing"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var taskID = primitive.NewObjectID()
var taskTitle = "Task Title"
var taskDescription = "Task Description"
var taskStatus = domain.TaskStatus("Started")
var taskDueDate = time.Now()
var updateTaskTitle = "Updated Task Title"

type TaskUseCaseTest struct {
	suite.Suite
	mockTaskRepo   *mocks.MockTaskRepository
	taskUsecase    domain.TaskUsecase
	testTask       *domain.Task
	taskUpdate     *domain.TaskUpdate
	updatedTask    *domain.Task
	nilTask        *domain.Task
	nilTaskUpdate  *domain.TaskUpdate
	invalidTask    *domain.Task
	contextTimeout time.Duration
}

func (suite *TaskUseCaseTest) SetupSuite() {
	suite.mockTaskRepo = new(mocks.MockTaskRepository)
	suite.taskUsecase = NewTaskUsecase(suite.mockTaskRepo, suite.contextTimeout)
	suite.contextTimeout = time.Second * 2
	suite.testTask = &domain.Task{
		ID:          taskID,
		Title:       taskTitle,
		Description: taskDescription,
		Status:      taskStatus,
		Due_Date:    taskDueDate,
	}
	suite.updatedTask = &domain.Task{
		ID:          taskID,
		Title:       updateTaskTitle,
		Description: taskDescription,
		Status:      taskStatus,
		Due_Date:    taskDueDate,
	}
	suite.taskUpdate = &domain.TaskUpdate{
		Title:       &updateTaskTitle,
		Description: &taskDescription,
		Status:      &taskStatus,
	}
	suite.nilTask = &domain.Task{}
	suite.invalidTask = &domain.Task{
		Title:       taskTitle,
		Description: taskDescription,
	}
	suite.nilTaskUpdate = &domain.TaskUpdate{}
}

// TestCreateTask tests the CreateTask method in the use case
func (suite *TaskUseCaseTest) TestCreateTask() {
	suite.mockTaskRepo.On("CreateTask", mock.Anything, suite.testTask).Return(nil)
	err := suite.taskUsecase.CreateTask(context.Background(), suite.testTask)
	suite.NoError(err)
}

// TestCreateTaskError tests the CreateTask method in the use case with an empty task
func (suite *TaskUseCaseTest) TestCreateTaskError() {
	suite.mockTaskRepo.On("CreateTask", mock.Anything, suite.nilTask).Return(assert.AnError)
	err := suite.taskUsecase.CreateTask(context.Background(), suite.nilTask)
	suite.Error(err)
}

// TestCreateInvalidTaskErrpr tests the CreateTask method in the use case with an invalid task
func (suite *TaskUseCaseTest) TestCreateInvalidTaskError() {
	suite.mockTaskRepo.On("CreateTask", mock.Anything, suite.invalidTask).Return(assert.AnError)
	err := suite.taskUsecase.CreateTask(context.Background(), suite.invalidTask)
	suite.Error(err)
}

// TestFetchTasks tests the FetchTasks method in the use case
func (suite *TaskUseCaseTest) TestFetchTasks() {
	suite.mockTaskRepo.On("Fetch", mock.Anything).Return([]domain.Task{}, nil)
	_, err := suite.taskUsecase.Fetch(context.Background())
	suite.NoError(err)
}

// TestUpdateTask tests the UpdateTask method in the use case
func (suite *TaskUseCaseTest) TestUpdateTask() {
	suite.mockTaskRepo.On("UpdateTask", mock.Anything, suite.taskUpdate, taskID.Hex()).Return(*suite.updatedTask, nil)
	updatedTask, err := suite.taskUsecase.UpdateTask(context.Background(), suite.taskUpdate, taskID.Hex())
	suite.Assert().Equal(updatedTask.Title, updateTaskTitle)
	suite.NoError(err)
}

// TestUpdateTaskError tests the UpdateTask method in the use case with an empty task
func (suite *TaskUseCaseTest) TestUpdateTaskError() {
	suite.mockTaskRepo.On("UpdateTask", mock.Anything, suite.nilTaskUpdate, taskID.Hex()).Return(*suite.nilTask, assert.AnError)
	_, err := suite.taskUsecase.UpdateTask(context.Background(), suite.nilTaskUpdate, taskID.Hex())
	suite.Error(err)
}

// TestDeleteTask tests the DeleteTask method in the use case
func (suite *TaskUseCaseTest) TestDeleteTask() {
	suite.mockTaskRepo.On("DeleteTask", mock.Anything, taskID.Hex()).Return(nil)
	err := suite.taskUsecase.DeleteTask(context.Background(), taskID.Hex())
	suite.NoError(err)
}

// TearDownSuite clears the mock
func (suite *TaskUseCaseTest) TearDownSuite() {
	suite.mockTaskRepo = nil
	suite.taskUsecase = nil
	suite.testTask = nil
}

func TestTaskUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseTest))
}
