// pkg/api/nextdate.go
package api

import (
	"log"
	"net/http"
	"time"

	"myfinproject/pkg/nextdate" // импортируем нашу функцию NextDate
)

// NextDateHandler обрабатывает запросы к /api/nextdate
// Формат запроса: /api/nextdate?now=20240126&date=20240126&repeat=y
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateStr == "" {
		http.Error(w, "Не указан параметр date", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "Не указан параметр repeat", http.StatusBadRequest)
		return
	}

	// Определяем текущую дату (now)
	var now time.Time
	if nowStr == "" {
		// Если now не указан, используем текущую дату
		now = time.Now()
		// Обрезаем время до начала дня
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		// Парсим переданную дату now
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат параметра now. Ожидается YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	// Вызываем функцию NextDate из пакета nextdate
	next, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		// В случае ошибки возвращаем текст ошибки
		log.Printf("Ошибка при вычислении следующей даты: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Успешно возвращаем дату
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
