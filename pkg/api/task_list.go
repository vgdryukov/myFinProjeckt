// pkg/api/task_list.go
// Получение списка задач
package api

import (
	"myfinproject/pkg/db"
	"net/http"
)

// TasksResp структура ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Функция tasksHandler обрабатывает GET-запросы на /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка: является ли метод запроса GET-запросом
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получение значения параметра search из запроса
	search := r.URL.Query().Get("search")

	// Получение задач из базы данных (максимум 50)
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, "Ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращение запрошенных задач
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})

}
