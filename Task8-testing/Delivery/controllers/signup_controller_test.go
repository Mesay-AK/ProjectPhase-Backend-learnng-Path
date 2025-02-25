package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	infrastructure "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Infrastructure"
	"github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type singupController struct {
	suite.Suite
	usecase       *mocks.SignupUsecase
	controller    *SignupController
	testingServer *httptest.Server
	timeout       time.Duration
	nilUser       domain.User
}

func (suite *singupController) SetupSuite() {
	suite.usecase = new(mocks.SignupUsecase)
	suite.controller = &SignupController{
		SignupUsecase: suite.usecase,
		Env:           bootstrap.NewEnv(3),
	}
	suite.timeout = time.Duration(suite.timeout) * time.Second
	suite.nilUser = domain.User{}
	router := gin.Default()
	router.POST("/user/register", suite.controller.Signup)

	// create and run the testing server
	testingServer := httptest.NewServer(router)
	suite.testingServer = testingServer

}

func (suite *singupController) TearDownSuite() {
	defer suite.testingServer.Close()
}

// TestSignup tests the Signup method in the controller
func (suite *singupController) TestSignup() {
	// Mock request payload
	requestPayload := domain.SignupRequest{
		First_Name: "John",
		Last_Name:  "Doe",
		Email:      "john.doe@example.com",
		User_Name:  "johndoe",
		Password:   "password123",
	}
	userRole := "USER"
	encryptedPassword, _ := infrastructure.HashPassword(requestPayload.Password)
	// Convert SignupRequest to User
	expectedUser := &domain.User{
		First_Name: &requestPayload.First_Name,
		Last_Name:  &requestPayload.Last_Name,
		Email:      &requestPayload.Email,
		User_Name:  &requestPayload.User_Name,
		Password:   &encryptedPassword,
		User_Role:  &userRole,
	}

	// Set up mock expectations
	suite.usecase.On("GetUserByEmail", mock.Anything, "john.doe@example.com").Return(suite.nilUser, errors.New("user not found"))
	suite.usecase.On("GetUserByUsername", mock.Anything, "johndoe").Return(suite.nilUser, errors.New("user not found"))
	suite.usecase.On("Create", mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(user *domain.User) bool {
		return *user.User_Name == *expectedUser.User_Name
	})).Return(nil)

	// Convert payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	assert.NoError(suite.T(), err)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/user/register", bytes.NewBuffer(requestBody))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Log the response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(suite.T(), err)

	// Assert the response body
	// Extract and compare message
	data := responseBody["data"].(map[string]interface{})
	message := data["message"].(string)
	suite.Equal("User created successfully", message)

	// Verify that the mock expectations were met
	suite.usecase.AssertExpectations(suite.T())
}

func TestSignupController(t *testing.T) {
	suite.Run(t, new(singupController))
}
