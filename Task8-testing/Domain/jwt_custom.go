package domain

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
