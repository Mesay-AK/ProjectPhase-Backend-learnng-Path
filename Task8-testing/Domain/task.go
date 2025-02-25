package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	NotStarted TaskStatus = "Not Started"
	Started    TaskStatus = "Started"
	Completed  TaskStatus = "Completed"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Status      TaskStatus         `json:"status" bson:"status" validate:"required,eq=Not Started|eq=Started|eq=Completed"`
	Due_Date    time.Time          `json:"due_date" bson:"due_date"`
}

type TaskUpdate struct {
	Title       *string     `json:"title" bson:"title"`
	Description *string     `json:"description" bson:"description"`
	Status      *TaskStatus `json:"status" bson:"status" validate:"omitempty,eq=Not Started|eq=Started|eq=Completed"`
	Due_Date    *time.Time  `json:"due_date" bson:"due_date"`
}

type TaskRepository interface {
	CreateTask(c context.Context, task *Task) error
	Fetch(c context.Context) ([]Task, error)
	GetByID(c context.Context, id string) (Task, error)
	UpdateTask(c context.Context, task *TaskUpdate, id string) (Task, error)
	DeleteTask(c context.Context, id string) error
}

type TaskUsecase interface {
	CreateTask(c context.Context, task *Task) error
	Fetch(c context.Context) ([]Task, error)
	GetByID(c context.Context, id string) (Task, error)
	UpdateTask(c context.Context, task *TaskUpdate, id string) (Task, error)
	DeleteTask(c context.Context, id string) error
}
