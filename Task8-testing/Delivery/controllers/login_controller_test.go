package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type loginController struct {
	suite.Suite
	usecase       *mocks.LoginUsecase
	controller    *LoginController
	testingServer *httptest.Server
	timeout       time.Duration
	nilUser       domain.User
}

func (suite *loginController) SetupSuite() {
	suite.usecase = new(mocks.LoginUsecase)
	suite.controller = &LoginController{
		LoginUsecase: suite.usecase,
		Env:          bootstrap.NewEnv(3),
	}
	suite.timeout = time.Duration(suite.timeout) * time.Second
	suite.nilUser = domain.User{}
	router := gin.Default()
	router.POST("/user/login", suite.controller.Login)

	// create and run the testing server
	testingServer := httptest.NewServer(router)
	suite.testingServer = testingServer

}

func (suite *loginController) TearDownSuite() {
	defer suite.testingServer.Close()
}

func strPtr(s string) *string {
	return &s
}

// TestLogin tests the Login method in the controller
func (suite *loginController) TestLogin() {

	// Mock request payload
	requestPayload := domain.LoginRequest{
		UserName: "johndoe",
		Password: "password123",
	}
	hashedPassword, _ := infrastructure.HashPassword(requestPayload.Password)

	// Mock response payload
	responsePayload := domain.User{
		First_Name: strPtr("John"),
		Last_Name:  strPtr("Doe"),
		Email:      strPtr("john.doe@example.com"),
		User_Name:  strPtr("johndoe"),
		Password:   &hashedPassword,
		User_Role:  strPtr("USER"),
	}

	accessTokenSecret := suite.controller.Env.AccessTokenSecret
	accessTokenExpiryHour := suite.controller.Env.AccessTokenExpiryHour

	// Mock the usecase
	suite.usecase.On("Login", mock.Anything, requestPayload).Return(responsePayload, nil)
	suite.usecase.On("GetUserByUsername", mock.Anything, "johndoe").Return(responsePayload, nil)
	suite.usecase.On("ComparePasswords", requestPayload.Password, hashedPassword).Return(nil)
	suite.usecase.On("CreateAccessToken", &responsePayload, accessTokenSecret, accessTokenExpiryHour).Return("valid_token", nil).Once()

	/// Convert payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	assert.NoError(suite.T(), err)

	// Make the request
	req, err := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/user/login", bytes.NewBuffer(requestBody))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Log the response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Decode the response body
	var responseBody domain.User
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(suite.T(), err)

	// Log the response body for debugging
	fmt.Printf("Response Body: %s\n", string(bodyBytes))

}

// TestLogin_IncorrectPassword tests the Login method that have an incorrect password in the controller
func (suite *loginController) TestLogin_IncorrectPassword() {

	requestPayload := domain.LoginRequest{
		UserName: "johndoe",
		Password: "password",
	}

	// Mock response payload
	responsePayload := domain.User{
		First_Name: strPtr("John"),
		Last_Name:  strPtr("Doe"),
		Email:      strPtr("john.doe@example.com"),
		User_Name:  strPtr("johndoe"),
		Password:   strPtr("password123"),
		User_Role:  strPtr("USER"),
	}
	// Mock the usecase
	suite.usecase.On("Login", mock.Anything, requestPayload).Return(responsePayload, nil)
	suite.usecase.On("GetUserByUsername", mock.Anything, "johndoe").Return(responsePayload, nil)
	suite.usecase.On("ComparePasswords", requestPayload.Password, responsePayload.Password).Return(nil)

	/// Convert payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	assert.NoError(suite.T(), err)

	// Make the request
	req, err := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/user/login", bytes.NewBuffer(requestBody))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	// assert the error message of Invalid password
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(suite.T(), err)

}

// TestLoginInvalidRequest tests the Login method that have an invalid request in the controller
func (suite *loginController) TestLoginInvalidRequest() {

	// Mock request payload
	requestPayload := domain.LoginRequest{
		Password: "password123",
	}

	// Mock the usecase
	suite.usecase.On("Login", mock.Anything, requestPayload).Return(domain.User{}, nil)

	/// Convert payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	assert.NoError(suite.T(), err)

	// Make the request
	req, err := http.NewRequest(http.MethodPost, suite.testingServer.URL+"/user/login", bytes.NewBuffer(requestBody))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func TestLoginController(t *testing.T) {
	suite.Run(t, new(loginController))
}
