package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("my-secret-key")

func GenerateToken(
	userID int,
	username string,
) (string, error) {

	claims := jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp": time.Now().
			Add(24 * time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(SecretKey)
}

func GetUserIDFromToken(
	tokenString string,
) (int, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		},
	)

	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := int(claims["id"].(float64))

	return userID, nil
}
