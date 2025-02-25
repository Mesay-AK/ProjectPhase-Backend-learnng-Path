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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var firstName = "John"
var lastName = "Doe"
var userName = "John"
var userRole = "Admin"
var password = "password"
var email = "test@test.com"

type UserSignupUsecase struct {
	suite.Suite
	mockUserRepo   *mocks.MockUserRepository
	signupUsecase  domain.SignupUsecase
	testUser       *domain.User
	halfUser       *domain.User
	nilUser        *domain.User
	contextTimeout time.Duration
}

func (suite *UserSignupUsecase) SetupSuite() {
	suite.mockUserRepo = new(mocks.MockUserRepository)
	suite.contextTimeout = time.Second * 2
	suite.signupUsecase = NewSignupUsecase(suite.mockUserRepo, suite.contextTimeout)
	suite.testUser = &domain.User{
		ID:         primitive.NewObjectID(),
		First_Name: &firstName,
		Last_Name:  &lastName,
		User_Name:  &userName,
		User_Role:  &userRole,
		Password:   &password,
		Email:      &email,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}
	suite.nilUser = &domain.User{}
	suite.halfUser = &domain.User{
		First_Name: &firstName,
		Last_Name:  &lastName,
		User_Name:  &userName,
		User_Role:  &userRole,
		Password:   &password,
	}

}

// TestCreateUser tests the CreateUser method in the use case
func (suite *UserSignupUsecase) TestCreateUser() {
	suite.mockUserRepo.On("CreateUser", mock.Anything, suite.testUser).Return(nil)

	err := suite.signupUsecase.Create(context.Background(), suite.testUser)

	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestCreateUserError tests the CreateUser method in the use case with an error
func (suite *UserSignupUsecase) TestCreateUserError() {
	suite.mockUserRepo.On("CreateUser", mock.Anything, suite.nilUser).Return(assert.AnError)

	err := suite.signupUsecase.Create(context.Background(), suite.nilUser)

	assert.Error(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestCreateInvalidUserError tests the CreateUser method in the use case with an invalid user
func (suite *UserSignupUsecase) TestCreateInvalidUserError() {
	suite.mockUserRepo.On("CreateUser", mock.Anything, suite.halfUser).Return(assert.AnError)

	err := suite.signupUsecase.Create(context.Background(), suite.halfUser)

	assert.Error(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// Teardown for the suite
func (suite *UserSignupUsecase) TearDownSuite() {
	suite.mockUserRepo = nil
	suite.signupUsecase = nil
	suite.testUser = nil
	suite.contextTimeout = 0
}

func TestUserSignupUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserSignupUsecase))
}
