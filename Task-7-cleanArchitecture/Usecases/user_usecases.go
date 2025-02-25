package usecase

import (
	"context"
	"fmt"
	domain "task_manager_clean/Domain"
	infrastructure "task_manager_clean/Infrastructure"
	"task_manager_clean/Repositories"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserUseCase struct {
	UserRepo *repository.UserRepository
	Validate *validator.Validate
}

func NewUserUseCase(ur *repository.UserRepository, v *validator.Validate) *UserUseCase {
	return &UserUseCase{
		UserRepo: ur,
		Validate: v,
	}
}


func (uCase *UserUseCase) Register(ctx context.Context, user *domain.User) error {
	if err := uCase.Validate.Struct(user); err != nil {
		fmt.Print("error while validating")
		return err
	}

	email := *user.Email

	count, err := uCase.UserRepo.CheckExistingUser(email)
	if err != nil {
		fmt.Print("error in the checkExistinguser")
		return err
	}

	hashedPassword, err := infrastructure.HashPassword(*user.Password)
	if err != nil {
		fmt.Print("error in the hashing")
		return err
	}
	user.Password = &hashedPassword

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.UserType == nil {
        user.UserType = new(string)
    }

	if count == 0 {
		*user.UserType = "ADMIN"
	} else {
		*user.UserType = "USER"
	}

	token, refreshToken, err := infrastructure.GenerateAllTokens(email, *user.Name, *user.UserType, *user.UserID)
	if err != nil {
		fmt.Print("error in the generating token")
		return err
	}
	user.Token = &token
	user.RefreshToken = &refreshToken

	return uCase.UserRepo.Register(ctx, user)
}

func (uCase *UserUseCase) LogIn(ctx context.Context, email, password string) (domain.User, string, string, error) {
	user, token, refreshToken, err := uCase.UserRepo.LogIn(ctx, email, password)
	if err != nil {
		return user, "", "", err
	}
	return user, token, refreshToken, nil
}

func (uCase *UserUseCase) GetUsers(ctx context.Context) ([]domain.User, error) {
	return uCase.UserRepo.GetUsers(ctx)
}

func (uCase *UserUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := uCase.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (uCase *UserUseCase) PromoteUser(ctx context.Context, id string, promoted *domain.User) (*domain.User, error) {
	user, err := uCase.UserRepo.PromoteUser(ctx, id, promoted)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (uCase *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return uCase.UserRepo.DeleteUser(ctx, id)
}
