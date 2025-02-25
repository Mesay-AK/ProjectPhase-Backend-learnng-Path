package infrastructure

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a password string and returns the hashed version of the password.
// It uses bcrypt.GenerateFromPassword to generate the hash with the default cost.
// If an error occurs during the hashing process, it returns an empty string and an error.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash Password")
	}

	return string(hashedPassword), nil
}

// ComparePasswords compares a plain text password with a hashed password.
// It uses bcrypt.CompareHashAndPassword to perform the comparison.
// If the passwords match, it returns nil. Otherwise, it returns an error.
func ComparePasswords(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return err
	}

	return nil
}
