// pkg/api/task_delete.go
// Удаление задачи
package api

import (
	"myfinproject/pkg/db"
	"net/http"
)

// Функция deleteTaskHandler обрабатывает DELETE-запросы на /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка: является ли метод запроса DELETE-запросом
	if r.Method != http.MethodDelete {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получение ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Удаление задачи из БД
	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Возвращение пустой JSON в случае успешного удаления задачи из БД
	writeJSON(w, map[string]interface{}{})

}
