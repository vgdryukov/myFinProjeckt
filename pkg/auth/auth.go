// pkg/auth/auth.go
// Аутентификация и работа с JWT-токенами
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	lifeTokenConfig = 8
	authEnabled     bool
	secretKey       []byte
	password        string
	passwordHash    string
)

// Функция init загружает конфигурацию
func init() {
	password = os.Getenv("TODO_PASSWORD")
	authEnabled = password != ""

	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	secretKey = []byte(secret)

	if authEnabled {
		hash := sha256.Sum256([]byte(password))
		passwordHash = hex.EncodeToString(hash[:])
	}
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func GenerateToken(pass string) (string, error) {
	// Используем глобальные переменные
	if authEnabled && pass != password {
		return "", errors.New("неверный пароль")
	}

	expirationTime := time.Now().Add(time.Duration(lifeTokenConfig) * time.Hour)

	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (bool, error) {
	// Используем глобальную переменную
	if !authEnabled {
		return true, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Сравнение с предвычисленным хэшем
		if claims.PasswordHash != passwordHash {
			return false, errors.New("пароль был изменен")
		}
		return true, nil
	}

	return false, errors.New("невалидный токен")
}

func IsAuthEnabled() bool {
	return authEnabled
}
