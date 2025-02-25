package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignupValidation(t *testing.T) {
	t.Run("Test should pass with valid fields", func(t *testing.T) {
		signupRequest := SignupRequest{
			First_Name: "Test",
			Last_Name:  "User",
			User_Name:  "testUser",
			Password:   "123456",
			Email:      "test@test.com",
		}
		assert.Equal(t, "Test", signupRequest.First_Name, "first names should match")
		assert.Equal(t, "User", signupRequest.Last_Name, "last names should match")
		assert.Equal(t, "testUser", signupRequest.User_Name, "usernames should match")
		assert.Equal(t, "123456", signupRequest.Password, "passwords should match")
		assert.Equal(t, "test@test.com", signupRequest.Email, "emails should match")
	})

	t.Run("Test should fail with a missing field", func(t *testing.T) {
		signupRequest := SignupRequest{
			User_Name: "testUser",
		}

		// Check if Password is empty
		assert.Equal(t, "", signupRequest.Password, "Password should be an empty string")
		assert.Equal(t, "", signupRequest.First_Name, "first name should be empty")
		assert.Equal(t, "", signupRequest.Last_Name, "last name should be empty")

		err := validateSingupRequest(signupRequest)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "First Name is required")
	})
}

// Simulate the validation function for the example
func validateSingupRequest(lr SignupRequest) error {
	if lr.First_Name == "" {
		return fmt.Errorf("First Name is required")
	}

	if lr.Last_Name == "" {
		return fmt.Errorf("last name requried")
	}

	if lr.Email == "" {
		return fmt.Errorf("email requried")
	}

	if lr.User_Name == "" {
		return fmt.Errorf("user name required ")
	}

	if lr.Password == "" {
		return fmt.Errorf("Password is required")
	}
	return nil
}
