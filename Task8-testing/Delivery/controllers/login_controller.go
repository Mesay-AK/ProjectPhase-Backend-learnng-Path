package controllers

import (
	"net/http"

	bootstrap "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Bootstrap"
	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	infrastructure "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Infrastructure"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

// Login is a method of the LoginController struct that handles the login functionality.
// It receives a gin.Context object as a parameter and expects a JSON payload containing a username and password.
// It validates the request payload, checks if the user exists, verifies the password, and creates an access token if successful.
// If any error occurs during the process, it returns an appropriate error response.
// On success, it returns an HTTP 200 status code along with the generated access token.
func (lc *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := lc.LoginUsecase.GetUserByUsername(c, request.UserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	err = infrastructure.ComparePasswords(request.Password, *user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
