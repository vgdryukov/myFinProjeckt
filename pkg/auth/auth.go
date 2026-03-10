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

var lifeTokenConfig = 8 // Время жизни токена установлено 8 часов

// Claims структура данных JWT-токена
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// Функция GenerateToken создает JWT-токен на основе пароля
func GenerateToken(password string) (string, error) {

	// Получение секретного ключа из переменной окружения, или используем значение по умолчанию
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	// Создание хэша пароля
	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	// Устанавка времени окончания жизни токена
	expirationTime := time.Now().Add(time.Duration(lifeTokenConfig) * time.Hour)

	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписание токена
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Функция ValidateToken проверяет валидность JWT токена
func ValidateToken(tokenString string) (bool, error) {

	// Получение пароля из переменной окружения
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		// Если пароль не установлен, значит аутентификация не требуется
		return true, nil
	}

	// Получение секретного ключа из переменной окружения
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	// Создание хэша текущего пароля для сравнения
	currentHash := sha256.Sum256([]byte(password))
	currentHashStr := hex.EncodeToString(currentHash[:])

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка что методом подписи токена является HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	// Проверка claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Сравниваем хэш пароля из токена с текущим
		if claims.PasswordHash != currentHashStr {
			return false, errors.New("пароль был изменен")
		}
		return true, nil
	}

	return false, errors.New("невалидный токен")
}

// Функция IsAuthEnabled проверяет, включена ли аутентификация
func IsAuthEnabled() bool {

	return os.Getenv("TODO_PASSWORD") != ""
}
