package jwt

import (
	"fmt"
	"rest/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	secretKey = "ololo114ololo"
)

func NewToken(user *domain.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("can't extract clains")
	}

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJwtToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false
	}

	if token.Valid {
		return true
	} else {
		return false
	}
}

func DecodeJwtToken(tokenString string) string {
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if value, ok := claims["email"].(string); ok && token.Valid {
		return value
	}
	return ""
}
