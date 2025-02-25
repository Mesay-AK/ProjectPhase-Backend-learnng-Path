package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type promoteContoller struct {
	suite.Suite
	usecase       *mocks.UserUsecase
	controller    *PromoteController
	testingServer *httptest.Server
	timeout       time.Duration
	nilUser       domain.User
}

func (suite *promoteContoller) SetupSuite() {
	suite.usecase = new(mocks.UserUsecase)
	suite.controller = &PromoteController{
		UserUsecase: suite.usecase,
		Env:         bootstrap.NewEnv(3),
	}
	suite.timeout = time.Duration(suite.timeout) * time.Second
	suite.nilUser = domain.User{}

	router := gin.Default()
	router.PUT("/user/promote/:id")

	testingServer := httptest.NewServer(router)
	suite.testingServer = testingServer
}

func (suite *promoteContoller) TearDownSuite() {
	defer suite.testingServer.Close()
}

func (suite *promoteContoller) TestPromote() {
	mockID := primitive.NewObjectID()

	promotedUser := domain.User{
		ID:         mockID,
		First_Name: strPtr("John"),
		Last_Name:  strPtr("Doe"),
		Email:      strPtr("john.doe@example.com"),
		User_Name:  strPtr("johndoe"),
		Password:   strPtr("password123"),
		User_Role:  strPtr("Admin"),
	}
	userID := mockID
	suite.usecase.On("Promote", mock.Anything, userID).Return(promotedUser, nil).Once()

	// Make the request
	req, err := http.NewRequest(http.MethodPut, suite.testingServer.URL+"/user/promote/"+userID.Hex(), nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

}

func TestPromoteUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(promoteContoller))
}
