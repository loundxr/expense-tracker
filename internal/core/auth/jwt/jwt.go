package core_jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func NewToken(user domain.User, secret string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uid":   user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(ttl).Unix(),
			"iat":   time.Now().Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return tokenString, nil
}
