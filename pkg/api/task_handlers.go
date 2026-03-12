// pkg/api/task_handlers.go
// Операции с одной задачей
package api

import (
	"encoding/json"
	"myfinproject/pkg/db"
	"net/http"
)

// Функция getTaskHandler обрабатывает GET-запросы на /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Получение ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получение задачи из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Возвращение задачи
	writeJSON(w, task)

}

// Функция updateTaskHandler обрабатывает PUT-запросы на /api/task
// updateTaskHandler обрабатывает PUT-запросы на /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодирование JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка факта, что указан ID
	if task.ID == "" {
		writeError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательного поля title
	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверка и корректировка даты
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление задачи в базе данных
	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращение пустого JSON
	writeJSON(w, map[string]interface{}{})
}
