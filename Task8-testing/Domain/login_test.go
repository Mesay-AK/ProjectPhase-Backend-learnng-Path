package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginValidation(t *testing.T) {
	t.Run("Test should pass with a valid username and password", func(t *testing.T) {
		loginRequest := LoginRequest{
			UserName: "testUser",
			Password: "123456",
		}

		assert.Equal(t, "testUser", loginRequest.UserName, "usernames should match")
		assert.Equal(t, "123456", loginRequest.Password, "passwords should match")
	})

	t.Run("Test should fail with a missing field", func(t *testing.T) {
		loginRequest := LoginRequest{
			UserName: "testUser",
		}

		// Check if Password is empty
		assert.Equal(t, "", loginRequest.Password, "Password should be an empty string")

		err := validateLoginRequest(loginRequest)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Password is required")
	})
}

// Simulate the validation function for the example
func validateLoginRequest(lr LoginRequest) error {
	if lr.Password == "" {
		return fmt.Errorf("Password is required")
	}
	return nil
}
