package data

import (
	"context"
	"errors"
	"fmt"
	"task_manager_with_jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskService struct {
	TaskCollection *mongo.Collection
}

func NewTaskService(db *mongo.Database) *TaskService {
	collection := db.Collection("tasks")
	return &TaskService{TaskCollection: collection}
}

var ErrTaskNotFound = errors.New("task not found")

func (service *TaskService) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	collection := service.TaskCollection

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}


// Task service reciever method to get task from the database collection
func (service *TaskService) GetTaskByID(id string) (*models.Task, error) {
	var task models.Task
	collection := service.TaskCollection

	filter := bson.M{"_id": id}

	err := collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}



// These are requests for authorized admin

// Task service reciever method to update a task on the database collection by it's ID

func (service *TaskService) AddTask(task models.Task) error {
	collection := service.TaskCollection

	_, err := collection.InsertOne(context.Background(), task)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("task already exists")
	}
	return err
}

func (service *TaskService) UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
	collection := service.TaskCollection

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": updatedTask,
	}

	result := collection.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var task models.Task
	if err := result.Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (service *TaskService) DeleteTask(id string) error {
	collection := service.TaskCollection

	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}