package usecases

import (
	"context"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
)

type promoteUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// NewPromoteUsecase creates a new instance of the PromoteUsecase struct that implements the domain.UserUsecase interface.
// It takes a userRepository of type domain.UserRepository and a timeout of type time.Duration as parameters.
// Returns a pointer to the PromoteUsecase struct.
func NewPromoteUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &promoteUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

// Promote promotes a user with the given ID.
// It takes a context and the user ID as parameters.
// It returns the promoted user and an error, if any.
func (pu *promoteUsecase) Promote(c context.Context, id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()
	return pu.userRepository.Promote(ctx, id)
}
