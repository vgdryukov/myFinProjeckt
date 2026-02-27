// pkg/api/tasks.go
package api

import (
	"myfinproject/pkg/db"
	"net/http"
)

// TasksResp структура ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы на /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр search из запроса
	search := r.URL.Query().Get("search")

	// Получаем задачи из базы данных (максимум 50)
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, "Ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем список задач
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
