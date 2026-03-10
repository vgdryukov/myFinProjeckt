// pkg/api/task_add.go
package api

import (
	"encoding/json"
	"errors"
	"myfinproject/pkg/db"
	"myfinproject/pkg/nextdate"
	"net/http"
	"strconv"
	"time"
)

// addTaskHandler обрабатывает POST-запросы на /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка при добавлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(nextdate.DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты. Ожидается YYYYMMDD")
	}

	if task.Repeat != "" {
		// Для новой задачи с правилом повторения:
		// если дата в прошлом, ищем ближайшую будущую дату от СЕГОДНЯ
		if !nextdate.AfterNow(t, now) {
			next, err := nextdate.NextDate(now, today, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else {
		// Без правила: если дата в прошлом, ставим сегодня
		if !nextdate.AfterNow(t, now) {
			task.Date = today
		}
	}

	return nil
}
