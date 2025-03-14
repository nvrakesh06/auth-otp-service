package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"os"
)

func GenerateJWT(phone string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"phone": phone,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}