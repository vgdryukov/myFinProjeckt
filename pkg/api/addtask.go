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

// Функция addTaskHandler обрабатывает POST-запросы на /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка: является ли метод запроса POST-запросом
	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка: обязательное поле title не пустое
	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверка и корректировка даты
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка при добавлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращение ID созданной задачи в случае ее успешного добавления
	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})

}

// Функция checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {

	now := time.Now()
	today := now.Format(nextdate.DateFormat)

	// Если дата не указана, ставим сегодняшнюю дату
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

		// Проверка правила повторения и получение следующей даты
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		// Если дата задачи меньше сегодняшней, используем вычисленную дату
		if !nextdate.AfterNow(t, now) {
			task.Date = next
		}

	} else {
		// Если правило не указано и дата меньше сегодняшней, ставим сегодняшнюю дату
		if !nextdate.AfterNow(t, now) {
			task.Date = today
		}
	}

	return nil
}
