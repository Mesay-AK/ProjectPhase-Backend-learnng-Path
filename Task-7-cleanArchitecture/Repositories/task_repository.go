// repository/task_repository_impl.go

package repository

import (
	"context"
	"errors"
	"fmt"
	domain "task_manager_clean/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	TaskCollection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	collection := db.Collection("tasks")
	return &TaskRepository{TaskCollection: collection}
}

var ErrTaskNotFound = errors.New("task not found")

func (repo *TaskRepository) GetTasks(ctx context.Context) ([]domain.Task, error) {
	var tasks []domain.Task

	cursor, err := repo.TaskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TaskRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	var task domain.Task

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format")
	}

	filter := bson.M{"_id": objectID}

	err = repo.TaskCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) AddTask(ctx context.Context, task *domain.Task) error {
	_, err := repo.TaskCollection.InsertOne(ctx, task)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("task already exists")
		}
		return err
	}
	return nil
}

func (repo *TaskRepository) UpdateTask(ctx context.Context, id string, updatedTask *domain.Task) (*domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format")
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": updatedTask,
	}

	result := repo.TaskCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var task domain.Task
	if err := result.Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid task ID format")
	}

	filter := bson.M{"_id": objectID}

	result, err := repo.TaskCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
