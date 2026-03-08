// pkg/api/taskdone.go
package api

import (
	"myfinproject/pkg/db"
	"myfinproject/pkg/nextdate"
	"net/http"
	"time"
)

// Функция taskDoneHandler обрабатывает POST-запросы на /api/task/done?id=<id>
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка: является ли метод запроса POST-запросом
	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных соответствующую полученному ID
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Проверка: есть ли правило повторения
	if task.Repeat == "" {
		// Если это одноразовая задача - удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {

		// Если это периодическая задача - вычисляем следующую дату
		now := time.Now()

		// Вычисление следующей даты выполнения задачи
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка при вычислении следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Обновление только даты задачи
		err = db.UpdateDate(id, next)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращение пустого JSON в случае успешного выполнения задачи
	writeJSON(w, map[string]interface{}{})

}
