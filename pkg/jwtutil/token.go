package jwtutil

import (
	"book-management-system/config"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	return []byte(config.GetEnv("JWT_SECRET", ""))
}

func GenerateToken(userID string, role string, liveTime time.Duration) (string, error) {
	expirationTime := time.Now().Add(liveTime)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ParseAndValidateToken(tokenString string) (*Claims, error) {
	fmt.Printf("JWT SECRET IS: %s", getJWTSecret())
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи токена")
		}

		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("неверный токен: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	return claims, nil
}
