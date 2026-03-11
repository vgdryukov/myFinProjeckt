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

	// Если дата не указана, ставим сегодняшнюю дату
	if task.Date == "" {
		task.Date = today
		return nil
	}

	// Парсим дату
	t, err := time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты. Ожидается YYYYMMDD")
	}

	// Если дата сегодня или в будущем - ничего не меняем!
	if nextdate.AfterNow(t, now) || task.Date == today {
		return nil
	}

	// Если дата в прошлом
	if task.Repeat == "" {
		// Без правила - ставим сегодня
		task.Date = today
	} else {
		// С правилом - вычисляем следующую дату от сегодня
		next, err := nextdate.NextDate(now, today, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
	}

	return nil
}
