package secure

import (
	"fmt"
	"time"

	"cmd/obyavigo/main.go/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT struct {
	secure *config.Secure
}

func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		secure: &cfg.Secure,
	}
}

func (j *JWT) GenerateJWTToken(userID *uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secure.SecretKey))
}

func (j *JWT) ParseJWTToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}
		return []byte(j.secure.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok || userIDStr == "" {
			return nil, fmt.Errorf("user_id not found in token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		return &userID, nil
	}
	return nil, fmt.Errorf("invalid token")
}
