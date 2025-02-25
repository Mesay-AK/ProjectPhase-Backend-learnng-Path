package domain

import (
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func (u *User) Validate() error {
	return validate.Struct(u)
}

func TestUserValidation(t *testing.T) {
	t.Run("Test should pass with a valid user", func(t *testing.T) {
		firstName := "Yohannes"
		lastName := "Solomon"
		email := "test@test.com"
		password := "123456"
		userRole := "USER"

		// a valid user with valid fields
		validUser := User{
			ID:         primitive.NewObjectID(),
			First_Name: &firstName,
			Last_Name:  &lastName,
			Email:      &email,
			Password:   &password,
			User_Role:  &userRole,
			Created_At: time.Now(),
			Updated_At: time.Now(),
		}

		assert.NotEmpty(t, validUser.ID, "ID should not be empty")
		assert.Equal(t, "Yohannes", *validUser.First_Name, "First Name should be Yohannes")
		assert.Equal(t, "Solomon", *validUser.Last_Name, "Last Name should be Solomon")
		assert.Equal(t, "test@test.com", *validUser.Email, "Email should be test@test.com")
		assert.Equal(t, "123456", *validUser.Password, "Password should be 123456")
		assert.Equal(t, "USER", *validUser.User_Role, "User Role should be USER")
		assert.NotEmpty(t, validUser.Created_At, "Created At should not be empty")
		assert.NotEmpty(t, validUser.Updated_At, "Updated At should not be empty")
	})

	t.Run("Test should fail with an invalid email format.", func(t *testing.T) {
		firstName := "Yohannes"
		lastName := "Solomon"
		email := "test"
		password := "123456"
		userRole := "USER"

		// an invalid user with invalid email format
		invalidUser := User{
			ID:         primitive.NewObjectID(),
			First_Name: &firstName,
			Last_Name:  &lastName,
			Email:      &email,
			Password:   &password,
			User_Role:  &userRole,
			Created_At: time.Now(),
			Updated_At: time.Now(),
		}
		err := invalidUser.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})

	t.Run("Test should fail with a short password.", func(t *testing.T) {
		firstName := "Yohannes"
		lastName := "Solomon"
		email := "test@test.com"
		password := "12"
		userRole := "USER"

		// an invalid user with invalid email format
		invalidUser := User{
			ID:         primitive.NewObjectID(),
			First_Name: &firstName,
			Last_Name:  &lastName,
			Email:      &email,
			Password:   &password,
			User_Role:  &userRole,
			Created_At: time.Now(),
			Updated_At: time.Now(),
		}
		err := invalidUser.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Password")
	})

	t.Run("Test should fail with a missing requried field", func(t *testing.T) {
		firstName := "Yohannes"
		lastName := "Solomon"
		email := "test@test.com"
		userRole := "USER"

		// an invalid user with invalid email format
		invalidUser := User{
			ID:         primitive.NewObjectID(),
			First_Name: &firstName,
			Last_Name:  &lastName,
			Email:      &email,
			User_Role:  &userRole,
			Created_At: time.Now(),
			Updated_At: time.Now(),
		}
		err := invalidUser.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Password")
	})

}
