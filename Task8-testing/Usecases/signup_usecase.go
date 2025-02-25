package usecases

import (
	"context"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	infrastructure "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Infrastructure"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// NewSignupUsecase creates a new instance of the SignupUsecase interface.
// It takes a userRepository of type domain.UserRepository and a timeout of type time.Duration as parameters.
// It returns a pointer to a signupUsecase struct that implements the SignupUsecase interface.
// The userRepository parameter is used to interact with the user data storage.
// The timeout parameter is used to set the context timeout for the usecase.
func NewSignupUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

// Create creates a new user.
// It takes a context.Context and a *domain.User as parameters.
// It returns an error.
// The function creates a new user by calling the CreateUser method of the userRepository.
func (su *signupUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	err := su.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail retrieves a user by their email address.
// It takes a context.Context and the email string as parameters.
// It returns the user domain.User and an error if any occurred.
func (su *signupUsecase) GetUserByEmail(c context.Context, email string) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByEmail(ctx, email)
}

// GetUserByUsername retrieves a user by their username.
// It takes a context.Context and the username as parameters.
// It returns the user domain object and an error if any.
func (su *signupUsecase) GetUserByUsername(c context.Context, userName string) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByUsername(ctx, userName)
}

// CreateAccessToken generates an access token for the given user with the specified secret and expiry.
// It returns the generated access token and any error encountered during the process.
func (su *signupUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return infrastructure.CreateAccessToken(user, secret, expiry)
}
