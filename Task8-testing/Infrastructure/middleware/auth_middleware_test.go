package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	infrastructure "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockAuthInfrastructure is a mock implementation of the infrastructure used in AuthMiddleware.
type MockAuthInfrastructure struct {
	mock.Mock
}

var firstName = "John"
var lastName = "Doe"
var userName = "johndoe"
var email = "johndoe@example.com"
var password = "password123"
var userRole = "USER"
var mockUser = domain.User{
	ID:         primitive.NewObjectID(),
	First_Name: &firstName,
	Last_Name:  &lastName,
	User_Name:  &userName,
	Email:      &email,
	Password:   &password,
	User_Role:  &userRole,
	Created_At: time.Now(),
	Updated_At: time.Now(),
}

func (m *MockAuthInfrastructure) IsAuthorized(token string, secret string) error {
	args := m.Called(token, secret)
	return args.Error(0)
}

func (m *MockAuthInfrastructure) CheckTokenExpiry(token string, secret string) (bool, error) {
	args := m.Called(token, secret)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthInfrastructure) ExtractClaims(token string, secret string) (map[string]interface{}, error) {
	args := m.Called(token, secret)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

type AuthMiddlewareSuite struct {
	suite.Suite
	env       *bootstrap.Env
	authInfra *MockAuthInfrastructure
}

func (suite *AuthMiddlewareSuite) SetupSuite() {
	suite.env = bootstrap.NewEnv(3)
	suite.authInfra = new(MockAuthInfrastructure)
}

func (suite *AuthMiddlewareSuite) TearDownSuite() {}

func (suite *AuthMiddlewareSuite) TestMissingAuthorizationHeader() {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(suite.env))
	router.POST("/tasks", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("POST", "/task", nil)
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"Authorization header is required"}`, w.Body.String())
}

// TestInvalidAuthorizationHeader tests the case where the Authorization header is invalid.
func (suite *AuthMiddlewareSuite) TestInvalidAuthorizationHeader() {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(suite.env))
	router.POST("/tasks", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("POST", "/task", nil)
	req.Header.Set("Authorization", "InvalidHeader")
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"Invalid authorization header"}`, w.Body.String())
}

// TestInvalidToken tests the case where the JWT token is invalid.
func (suite *AuthMiddlewareSuite) TestInvalidToken() {
	// Mock the IsAuthorized function to return an error
	mockAuth := new(MockAuthInfrastructure)
	mockAuth.On("IsAuthorized", mock.Anything, mock.Anything).Return(fmt.Errorf("token contains an invalid number of segments"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AuthMiddleware(suite.env))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer InvalidToken")
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"token contains an invalid number of segments"}`, w.Body.String())
}

func (suite *AuthMiddlewareSuite) TestValidToken() {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	token, _ := infrastructure.CreateAccessToken(&mockUser, string(suite.env.AccessTokenSecret), suite.env.AccessTokenExpiryHour)
	router.Use(AuthMiddleware(suite.env))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Run the test suite
func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}
