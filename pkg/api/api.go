// pkg/api/api.go
// Регитрация маршрутов и общие функции API
package api

import (
	"encoding/json"
	"net/http"
)

// Функция Init регистрирует все API обработчики
func Init() {

	// Публичные маршруты (не требующие аутентификации)
	http.HandleFunc("/api/signin", signinHandler)

	// Защищенные маршруты (требуют аутентификации)
	http.HandleFunc("/api/nextdate", AuthMiddleware(NextDateHandler))
	http.HandleFunc("/api/task", AuthMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(taskDoneHandler))

}

// Функция taskHandler распределяет запросы по методам для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

}

// Функция writeJSON отправляет данные в формате JSON
func writeJSON(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
	}

}

// Функция writeError отправляет сообщение об ошибке в формате JSON
func writeError(w http.ResponseWriter, message string, statusCode int) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})

}
