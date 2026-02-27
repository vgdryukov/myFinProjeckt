// pkg/api/taskdone.go
package api

import (
	"myfinproject/pkg/db"
	"myfinproject/pkg/nextdate"
	"net/http"
	"time"
)

// taskDoneHandler обрабатывает POST-запросы на /api/task/done?id=<id>
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
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

	// Получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Проверяем, есть ли правило повторения
	if task.Repeat == "" {
		// Одноразовая задача - удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Периодическая задача - вычисляем следующую дату
		now := time.Now()

		// Вычисляем следующую дату выполнения
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка при вычислении следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем только дату задачи
		err = db.UpdateDate(id, next)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{})
}
