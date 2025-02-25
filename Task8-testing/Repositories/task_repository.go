package repositories

import (
	"errors"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

// NewTaskRepository creates a new instance of the TaskRepository interface.
// It takes a mongo.Database and a collection name as parameters and returns a pointer to a taskRepository struct.
// The taskRepository struct implements the TaskRepository interface and provides methods for interacting with the database.
func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

// CreateTask creates a new task in the repository.
// It takes a context.Context and a *domain.Task as parameters.
// The task is inserted into the collection specified in the taskRepository struct.
// Returns an error if there was a problem inserting the task.
func (tr *taskRepository) CreateTask(c context.Context, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.InsertOne(c, task)
	if err != nil {
		return err
	}

	return nil
}

// Fetch retrieves all tasks from the database.
// It returns a slice of domain.Task and an error if any.
func (tr *taskRepository) Fetch(c context.Context) ([]domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	var tasks []domain.Task
	err = cursor.All(c, &tasks)
	if err != nil {
		return []domain.Task{}, err
	}

	return tasks, err
}

// GetByID retrieves a task from the database by its ID.
// It takes a context.Context and the ID of the task as parameters.
// It returns the retrieved task and an error, if any.
func (tr *taskRepository) GetByID(c context.Context, id string) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	var task domain.Task
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return task, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&task)
	return task, err
}

// UpdateTask updates a task in the repository with the given ID.
// It takes a context, a TaskUpdate object containing the updated task fields, and the ID of the task to update.
// It returns the updated task and an error if any.
func (tr *taskRepository) UpdateTask(c context.Context, task *domain.TaskUpdate, id string) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	updateFields := make(bson.M)
	if task.Title != nil {
		updateFields["title"] = task.Title
	}

	if task.Status != nil {
		updateFields["status"] = task.Status
	}

	if task.Description != nil {
		updateFields["description"] = task.Description
	}

	var taskResponse domain.Task
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return taskResponse, err
	}

	result, err := collection.UpdateOne(
		c,
		bson.M{"_id": idHex},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		return taskResponse, err
	}

	if result.MatchedCount == 0 {
		return domain.Task{}, errors.New("no task found with the given ID")
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&taskResponse)
	if err != nil {
		return domain.Task{}, err
	}

	return taskResponse, nil
}

// DeleteTask deletes a task from the database.
// It takes a context.Context and the ID of the task to be deleted as parameters.
// It returns an error if the deletion fails.
func (tr *taskRepository) DeleteTask(c context.Context, id string) error {
	collection := tr.database.Collection(tr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	response, err := collection.DeleteOne(c, bson.M{"_id": idHex})
	if err != nil {
		return err
	}

	if response.DeletedCount == 0 {
		return errors.New("unable to delete Task")
	}

	return nil

}
