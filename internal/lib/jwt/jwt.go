package jwt

import (
	"fmt"
	"log"
	"rest/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	secretKey = "ololo114ololo"
)

type jwtClaims struct {
	jwt.Claims
}

func NewToken(user *domain.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("can't extract claims")
	}

	claims["uid"] = int(user.ID)
	// claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	log.Println("[NewToken]", claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*domain.User, bool) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, false
	}

	if token.Valid {
		var user domain.User
		if uid, ok := claims["uid"].(float64); ok {
			user.ID = int64(uid)
		}
		return &user, true
	} else {
		return nil, false
	}
}
