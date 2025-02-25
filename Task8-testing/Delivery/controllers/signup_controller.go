package controllers

import (
	"net/http"
	"time"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	infrastructure "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Infrastructure"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *bootstrap.Env
}

// Signup handles the signup request and creates a new user.
// It expects a JSON payload containing the user's signup details.
// If the request is valid and the user does not already exist, it creates a new user and returns a success response.
// If the request is invalid or the user already exists, it returns an appropriate error response.
func (sc *SignupController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sc.SignupUsecase.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	_, err = sc.SignupUsecase.GetUserByUsername(c, request.User_Name)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	encryptedPassword, err := infrastructure.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	request.Password = encryptedPassword
	userRole := "USER"

	user := domain.User{
		ID:         primitive.NewObjectID(),
		First_Name: &request.First_Name,
		Last_Name:  &request.Last_Name,
		Email:      &request.Email,
		User_Name:  &request.User_Name,
		Password:   &request.Password,
		User_Role:  &userRole,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	err = sc.SignupUsecase.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := infrastructure.CreateAccessToken(&user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	singupResponse := domain.SignupResponse{
		Message:      "User created successfully",
		Access_Token: accessToken,
		User_ID:      user.ID.Hex(),
	}

	c.JSON(http.StatusCreated, gin.H{"data": singupResponse})

}
