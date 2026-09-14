package core_jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UserClaims struct {
	ID   int64
	Role string
}

func NewToken(user domain.User, secret string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uid":   user.ID,
			"email": user.Email,
			"role":  user.Role,
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

func ParseToken(tokenString string, secret string) (UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return UserClaims{}, fmt.Errorf("parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uidFloat, ok := claims["uid"].(float64)
		if !ok {
			return UserClaims{}, fmt.Errorf("uid claim not found in token")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return UserClaims{}, fmt.Errorf("role claim not found in token")
		}
		return UserClaims{
			ID:   int64(uidFloat),
			Role: role,
		}, nil
	}

	return UserClaims{}, fmt.Errorf("invalid token")
}
