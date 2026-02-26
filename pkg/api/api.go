// pkg/api/api.go
package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// Init регистрирует все API обработчики
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
}

// taskHandler распределяет запросы по методам
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// writeJSON отправляет данные в формате JSON
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
	}
}

// writeError отправляет сообщение об ошибке в формате JSON
func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// afterNow проверяет, что дата date больше даты now
func afterNow(date, now time.Time) bool {
	// Нормализуем даты до начала дня
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}
