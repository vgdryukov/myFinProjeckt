// pkg/api/deletetask.go
package api

import (
	"myfinproject/pkg/db"
	"net/http"
)

// deleteTaskHandler обрабатывает DELETE-запросы на /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodDelete {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из базы данных
	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{})
}
