// pkg/api/signin.go
// Аутентификация пользователя
package api

import (
	"encoding/json"
	"myfinproject/pkg/auth"
	"net/http"
	"os"
)

// SigninRequest - структура запроса на аутентификацию
type SigninRequest struct {
	Password string `json:"password"`
}

// SigninResponse - структура успешного ответа
type SigninResponse struct {
	Token string `json:"token"`
}

// Функция signinHandler обрабатывает POST-запросы на /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка: является ли метод запроса POST-запросом
	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получение пароля из переменной окружения
	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		// Если пароль не установлен, аутентификация не требуется
		writeError(w, "Аутентификация не настроена", http.StatusBadRequest)
		return
	}

	// Декодирование JSON из тела запроса
	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка пароля
	if req.Password != expectedPassword {
		writeError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	// Генерация JWT-токена
	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeError(w, "Ошибка при создании токена: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращение токена
	writeJSON(w, SigninResponse{Token: token})

}
