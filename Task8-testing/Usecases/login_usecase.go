package usecases

import (
	"context"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	infrastructure "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Infrastructure"
)

type loginUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// NewLoginUsecase creates a new instance of the LoginUsecase implementation.
// It takes a userRepository of type domain.UserRepository and a timeout of type time.Duration as parameters.
// It returns a domain.LoginUsecase interface.
// The userRepository is responsible for accessing and manipulating user data.
// The timeout specifies the maximum duration for the usecase operations.
func NewLoginUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

// GetUserByEmail retrieves a user by their email address.
// It takes a context.Context and the email string as parameters.
// It returns the user domain.User and an error if any.
func (lu *loginUsecase) GetUserByEmail(c context.Context, email string) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

// GetUserByUsername retrieves a user by their username.
// It takes a context.Context and a userName string as parameters.
// It returns a domain.User and an error.
// The context.Context is used for managing the execution deadline and cancellation.
// The userName parameter specifies the username of the user to retrieve.
// The function returns the retrieved user and any error encountered during the retrieval process.
func (lu *loginUsecase) GetUserByUsername(c context.Context, userName string) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByUsername(ctx, userName)
}

// CreateAccessToken generates an access token for the given user with the specified secret and expiry.
// It returns the generated access token and any error encountered during the process.
func (lu *loginUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return infrastructure.CreateAccessToken(user, secret, expiry)
}
