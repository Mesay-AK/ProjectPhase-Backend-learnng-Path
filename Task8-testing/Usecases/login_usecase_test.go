package usecases

import (
	"context"
	"testing"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserLoginUsecase struct {
	suite.Suite
	mockUserRepo   *mocks.MockUserRepository
	loginUsecase   domain.LoginUsecase
	testUser       *domain.LoginRequest
	invalidUser    *domain.LoginRequest
	nilUser        *domain.LoginRequest
	contextTimeout time.Duration
}

func (suite *UserLoginUsecase) SetupSuite() {
	suite.mockUserRepo = new(mocks.MockUserRepository)
	suite.contextTimeout = time.Second * 2
	suite.loginUsecase = NewLoginUsecase(suite.mockUserRepo, suite.contextTimeout)
	suite.testUser = &domain.LoginRequest{
		UserName: "Johna",
		Password: "1234566",
	}
	suite.nilUser = &domain.LoginRequest{}
	suite.invalidUser = &domain.LoginRequest{
		UserName: "Johna"}
}

// TestLoginUser tests the LoginUser method in the use case
func (suite *UserLoginUsecase) TestLoginUser() {
	ctx, cancel := context.WithTimeout(context.Background(), suite.contextTimeout)
	defer cancel()
	suite.mockUserRepo.On("GetByUsername", mock.Anything, suite.testUser.UserName).Return(domain.User{}, nil)
	_, err := suite.loginUsecase.GetUserByUsername(ctx, suite.testUser.UserName)
	assert.NoError(suite.T(), err)
}

// TestLoginUserError tests the loginuserUsecase method with an empty user
func (suite *UserLoginUsecase) TestLoginUserError() {
	ctx, cancel := context.WithTimeout(context.Background(), suite.contextTimeout)
	defer cancel()
	suite.mockUserRepo.On("GetByUsername", mock.Anything, suite.nilUser.UserName).Return(domain.User{}, assert.AnError)
	_, err := suite.loginUsecase.GetUserByUsername(ctx, suite.nilUser.UserName)
	assert.Error(suite.T(), err)
}

// TeaeardownSuite clears the mock
func (suite *UserLoginUsecase) TearDownSuite() {
	suite.mockUserRepo = nil
	suite.loginUsecase = nil
	suite.testUser = nil
}

func TestUserLoginUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserLoginUsecase))
}
