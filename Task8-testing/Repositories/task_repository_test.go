package repositories

import (
	"context"
	"testing"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskId = primitive.NewObjectID()
var title = "New Task"
var description = "This is a new task"
var status = "Started"
var due_date = time.Now().UTC().Truncate(24 * time.Hour)

type TaskRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection string
	repo       domain.TaskRepository
}

func (suite *TaskRepositoryTestSuite) SetupSuite() {
	// Connect to the MongoDB test instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err)

	// Ping the MongoDB instance to ensure connection is established
	err = client.Ping(context.Background(), nil)
	suite.Require().NoError(err)

	suite.client = client
	suite.db = client.Database("testTaskManager")
	suite.collection = "tasks"
	suite.repo = NewTaskRepository(*suite.db, suite.collection)
}

func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	// Drop the test database after all tests have run
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	suite.Require().NoError(err)
}

func (suite *TaskRepositoryTestSuite) TestCreateTask() {
	task := domain.Task{
		ID:          taskId,
		Title:       title,
		Description: description,
		Status:      domain.TaskStatus(status),
		Due_Date:    due_date,
	}

	err := suite.repo.CreateTask(context.Background(), &task)
	suite.Require().NoError(err)

	// Verify task wask inserted
	collection := suite.db.Collection(suite.collection)
	var result domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&result)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), task.Title, result.Title)
	suite.NoError(err, "no error when creating and finding the same task")
}

func (suite *TaskRepositoryTestSuite) TestFetch() {
	fetchedTasks, err := suite.repo.Fetch(context.Background())
	suite.Require().NoError(err)

	// Verify the correct number of tasks were fetched
	assert.Equal(suite.T(), 1, len(fetchedTasks))

	for _, tasks := range fetchedTasks {
		assert.Equal(suite.T(), title, tasks.Title, "Same task title  with the previously created task")
		assert.Equal(suite.T(), description, tasks.Description, "Same task description with the previously created task")
		assert.Equal(suite.T(), status, string(tasks.Status), "Same task status with the previously created task")
		assert.Equal(suite.T(), due_date, tasks.Due_Date, "Same task due_date with the previously created task")
	}

}

func (suite *TaskRepositoryTestSuite) TestGetByID() {
	task, err := suite.repo.GetByID(context.Background(), taskId.Hex())
	suite.Require().NoError(err)

	assert.Equal(suite.T(), title, task.Title, "Same task title with the previously created task")
	assert.Equal(suite.T(), description, task.Description, "Same task description with the previously created task")
	assert.Equal(suite.T(), status, string(task.Status), "Same task status with the previously created task")
	assert.Equal(suite.T(), due_date, task.Due_Date, "Same task due_date with the previously created task")
}

func (suite *TaskRepositoryTestSuite) TestUpdateTask() {
	newTitle := "Updated Task"
	newDescription := "This is an updated task"
	newStatus := domain.TaskStatus("Completed")
	newDueDate := time.Now().UTC().Truncate(24 * time.Hour)

	taskUpdate := domain.TaskUpdate{
		Title:       &newTitle,
		Description: &newDescription,
		Status:      &newStatus,
		Due_Date:    &newDueDate,
	}

	_, err := suite.repo.UpdateTask(context.Background(), &taskUpdate, taskId.Hex())
	suite.Require().NoError(err)

	// find out if the task is updated
	collection := suite.db.Collection(suite.collection)
	var result domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": taskId}).Decode(&result)
	suite.Require().NoError(err)

	// Check if its the same task
	taskFound, err := suite.repo.GetByID(context.Background(), taskId.Hex())
	suite.Require().NoError(err)

	// check if the fields are updated
	assert.Equal(suite.T(), newTitle, taskFound.Title, "Same updated task title with the previously updated task")
	assert.Equal(suite.T(), newDescription, taskFound.Description, "Same updated task description with the previously updated task")
	assert.Equal(suite.T(), newStatus, taskFound.Status, "Same updated task status with the previously updated task")
	assert.Equal(suite.T(), newDueDate, taskFound.Due_Date, "Same updated task due_date with the previously updated task")
}

func (suite *TaskRepositoryTestSuite) TestDeleteTask() {
	newID := primitive.NewObjectID()

	// create a new task then remove it
	task := domain.Task{
		ID:          newID,
		Title:       "new task delete",
		Description: "This is a new task to be deleted",
		Status:      domain.TaskStatus("Not Started"),
		Due_Date:    time.Now().UTC().Truncate(24 * time.Hour),
	}

	err := suite.repo.CreateTask(context.Background(), &task)
	suite.Require().NoError(err)

	// Verify the task is inserted
	collection := suite.db.Collection(suite.collection)
	var result domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": newID}).Decode(&result)
	suite.Require().NoError(err)

	err = suite.repo.DeleteTask(context.Background(), newID.Hex())
	suite.Require().NoError(err)

	var newResult domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": result.ID}).Decode(&newResult)
	assert.Error(suite.T(), err, "task not found after deletion")
}

func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
