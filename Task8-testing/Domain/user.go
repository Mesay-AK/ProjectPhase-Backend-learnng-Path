package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_Name *string            `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	Last_Name  *string            `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	User_Name  *string            `json:"user_name" bson:"user_name" validate:"required,min=5"`
	Email      *string            `json:"email" bson:"email" validate:"required,email"`
	Password   *string            `json:"password" bson:"password" validate:"required,min=6"`
	User_Role  *string            `json:"user_role" bson:"user_role" validate:"omitempty,eq=Admin|eq=USER"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
	Updated_At time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserUpdate struct {
	First_Name *string `json:"first_name" bson:"first_name" validate:"omitempty,min=2,max=100"`
	Last_Name  *string `json:"last_name" bson:"last_name" validate:"omitempty,min=2,max=100"`
	User_Name  *string `json:"user_name" bson:"user_name" validate:"omitempty,min=5"`
	Email      *string `json:"email" bson:"email" validate:"omitempty,email"`
	Password   *string `json:"password" bson:"password" validate:"omitempty,min=6"`
	User_Role  *string `json:"user_role" bson:"user_role" validate:"omitempty,eq=Admin|eq=USER"`
}

type UserRepository interface {
	CreateUser(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByID(c context.Context, id string) (User, error)
	GetByUsername(c context.Context, userName string) (User, error)
	GetByEmail(c context.Context, email string) (User, error)
	UpdateUser(c context.Context, id string, user UserUpdate) (User, error)
	Promote(c context.Context, id string) (User, error)
}

type UserUsecase interface {
	Promote(c context.Context, id string) (User, error)
}
