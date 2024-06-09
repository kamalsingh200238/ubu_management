package utils

import (
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type CustomJwtClaims struct {
	jwt.RegisteredClaims
	PersonID int    `json:"id"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
}

type Role string

const (
	StudentRole   Role = "student"
	PresidentRole Role = "president"
	Coordinator   Role = "coordinator"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJwtToken(claims CustomJwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		slog.Error("invalid signature in token", err)
		return "", err
	}

	return tokenString, nil
}

func ParseAndValidateJWT(tokenString string) (*CustomJwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		slog.Error("cannot parse jwt token", err)
		return nil, err
	}

	// access the claims
	if claims, ok := token.Claims.(*CustomJwtClaims); ok && token.Valid {
		return claims, nil
	}

	slog.Error("token not valid", err)
	return nil, err
}
