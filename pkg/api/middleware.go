// pkg/api/middleware.go
// Аутентификация, проверка и валидация токена
package api

import (
	"myfinproject/pkg/auth"
	"net/http"
)

// Функция AuthMiddleware проверяет аутентификацию для защищенных маршрутов
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Проверка: включена ли аутентификация
		if auth.IsAuthEnabled() {

			// Получение токена из куков
			cookie, err := r.Cookie("token")
			if err != nil {

				// Если куков нет, проверяем заголовок Authorization (для API запросов)
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					writeError(w, "Authentification required", http.StatusUnauthorized)
					return
				}

				// Формат: "Bearer <token>"
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {

					// Валидация токена из заголовка
					valid, err := auth.ValidateToken(authHeader[7:])
					if err != nil || !valid {
						writeError(w, "Неверный токен", http.StatusUnauthorized)
						return
					}

				} else {
					writeError(w, "Неверный формат авторизации", http.StatusUnauthorized)
					return
				}

			} else {

				// Валидация токена из куков
				valid, err := auth.ValidateToken(cookie.Value)
				if err != nil || !valid {
					writeError(w, "Неверный токен", http.StatusUnauthorized)
					return
				}
			}
		}

		// Аутентификация пройдена или не требуется
		next(w, r)
	})

}
