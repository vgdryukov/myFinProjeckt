// pkg/api/addtask.go
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
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
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

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка при добавлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ с ID созданной задачи
	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(nextdate.DateFormat)

	// Если дата не указана, ставим сегодня
	if task.Date == "" {
		task.Date = today
	}

	// Парсим дату
	t, err := time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты. Ожидается YYYYMMDD")
	}

	// Если правило повторения указано, проверяем его
	if task.Repeat != "" {
		// Проверяем правило и получаем следующую дату
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		// Если дата задачи меньше сегодняшней, используем вычисленную
		if !afterNow(t, now) {
			task.Date = next
		}
	} else {
		// Если правило не указано и дата меньше сегодняшней, ставим сегодня
		if !afterNow(t, now) {
			task.Date = today
		}
	}

	return nil
}
