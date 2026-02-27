// pkg/api/task.go
package api

import (
	"encoding/json"
	"myfinproject/pkg/db"
	"net/http"
)

// getTaskHandler обрабатывает GET-запросы на /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Возвращаем задачу
	writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT-запросы на /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что указан ID
	if task.ID == "" {
		writeError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных
	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{})
}
