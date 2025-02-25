package Infrastructure

import (
    "golang.org/x/crypto/bcrypt"

)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func VerifyPassword(userPassword, foundPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(userPassword))
    check := true
    mssg := ""

    if err != nil {
        mssg = "Invalid email or password"
        check = false
    }

    return check, mssg
}

