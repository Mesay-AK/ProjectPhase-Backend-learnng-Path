package infrastructure

import (
	"fmt"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/golang-jwt/jwt"
)

// CreateAccessToken generates an access token for the given user with the specified secret and expiry time.
// It returns the access token as a string and any error encountered during the process.
func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := domain.JwtCustomClaims{
		UserID:   user.ID.Hex(),
		UserName: *user.User_Name,
		Role:     *user.User_Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, err
}

// IsAuthorized checks if the provided JWT token is authorized using the given secret.
// It parses the token and verifies the signing method using the HMAC algorithm.
// If the token is not valid or the signing method is unexpected, it returns an error.
func IsAuthorized(tokenString string, secret string) error {
	_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ExtractIDFromToken extracts the user ID from a JWT token.
// It takes the token string and the secret key as input parameters.
// It returns the extracted user ID as a string and an error if any.
func ExtractIDFromToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims["UserID"].(string), nil
}

// ExtractClaims extracts the claims from a JWT token.
// It takes the token string and secret as input parameters.
// It returns a map[string]interface{} containing the extracted claims and an error if any.
func ExtractClaims(tokenString string, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims := make(map[string]interface{})
	claims["UserID"] = token.Claims.(jwt.MapClaims)["user_id"].(string)
	claims["UserName"] = token.Claims.(jwt.MapClaims)["user_name"].(string)
	claims["Role"] = token.Claims.(jwt.MapClaims)["role"].(string)
	claims["exp"] = token.Claims.(jwt.MapClaims)["exp"].(float64)

	return claims, nil
}

// ExtractRoleFromToken extracts the role from the given JWT token.
// It takes the token string and secret key as input parameters.
// It returns the role as a string and an error if the token is invalid.
func ExtractRoleFromToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims["Role"].(string), nil
}

// CheckTokenExpiry checks the expiry of a JWT token.
// It takes the token string and the secret key as input parameters.
// It returns a boolean value indicating whether the token has expired or not, and an error if any.
func CheckTokenExpiry(tokenString string, secret string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return false, nil
		}
	}

	return true, nil
}
