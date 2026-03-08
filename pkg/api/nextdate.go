// pkg/api/nextdate.go
package api

import (
	"log"
	"myfinproject/pkg/nextdate"
	"net/http"
	"time"
)

// Функция NextDateHandler обрабатывает запросы к /api/nextdate
// Формат запроса: /api/nextdate?now=20240126&date=20240126&repeat=y
func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	// Получаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверка наличия в запросе обязательных параметров date и repeat
	if dateStr == "" {
		writeError(w, "Не указан параметр date", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		writeError(w, "Не указан параметр repeat", http.StatusBadRequest)
		return
	}

	// Определение текущей даты (now)
	now, err := parseNow(nowStr)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызов функции NextDate из пакета nextdate для вычисления следующей даты
	next, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		// В случае ошибки возвращаем текст ошибки
		log.Printf("Ошибка при вычислении следующей даты: %v", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращение следующей даты в случае её успешного определения
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))

}

// Функция parseNow парсит параметр now или возвращает текущую дату
func parseNow(nowStr string) (time.Time, error) {

	if nowStr == "" {
		// Если now не указан, используем текущую дату
		now := time.Now()
		// Обрезаем время до начала дня
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}

	// Парсим переданную дату now
	now, err := time.Parse(nextdate.DateFormat, nowStr)
	if err != nil {
		return time.Time{}, err
	}

	return now, nil
}
