package Domain

import (
	"context"
	"time"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          string    `bson:"_id" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	DueDate     time.Time `bson:"due_date" json:"due_date"`
	Status      string    `bson:"status" json:"status"`
}

type User struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          *string            `json:"name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Token         *string            `json:"token"`
	UserType      *string            `json:"user_type"`
	RefreshToken  *string            `json:"refresh_token"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	UserID        *string            `json:"user_id"`
}

type UserUsecase interface {
    Register(ctx context.Context, user *User) error
    LogIn(ctx context.Context, email, password string) (User, string, string, error)
    GetUsers(ctx context.Context) ([]User, error)
    GetUserByID(ctx context.Context, id string) (*User, error)
    PromoteUser(ctx context.Context, id string, user *User) (*User, error)
    DeleteUser(ctx context.Context, id string) error
}

type TaskUsecase interface {
    AddTask(ctx context.Context, task *Task) error
    GetTaskByID(ctx context.Context, id string) (*Task, error)
    GetTasks(ctx context.Context) ([]Task, error)
    UpdateTask(ctx context.Context, id string, task *Task) (*Task, error)
    DeleteTask(ctx context.Context, id string) error
}


var (
    ErrUserNotFound = errors.New("user not found")
    ErrTaskNotFound = errors.New("task not found")
    
)