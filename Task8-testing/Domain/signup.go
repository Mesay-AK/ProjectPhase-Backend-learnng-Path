package domain

import "context"

type SignupRequest struct {
	First_Name string `json:"first_name" validate:"required,min=2,max=100"`
	Last_Name  string `json:"last_name" validate:"required,min=2,max=100"`
	User_Name  string `json:"user_name" validate:"required,min=5"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
}

type SignupResponse struct {
	Message      string `json:"message"`
	User_ID      string `json:"user_id"`
	Access_Token string `json:"access_token"`
}

type SignupUsecase interface {
	Create(c context.Context, user *User) error
	GetUserByEmail(c context.Context, email string) (user User, err error)
	GetUserByUsername(c context.Context, userName string) (user User, err error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
}
