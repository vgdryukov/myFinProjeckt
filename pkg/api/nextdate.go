// pkg/api/nextDate.go
// API для вычисления следующей даты
package api

import (
	"errors"
	"log"
	"myfinproject/pkg/nextdate"
	"net/http"
	"time"
)

const (
	maxDateLen   = 8 // YYYYMMDD
	maxRepeatLen = 100
)

// Функция NextDateHandler обрабатывает запросы к /api/nextdate
// Формат запроса: /api/nextdate?now=20240126&date=20240126&repeat=y
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Защита от паники
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Panic в NextDateHandler: %v", rec)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}()

	// Проверка: метод запроса является GET
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	// Получение параметров из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверка: наличии обязательных параметров date, repeat
	if dateStr == "" {
		writeError(w, "Не указан параметр date", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		writeError(w, "Не указан параметр repeat", http.StatusBadRequest)
		return
	}

	// Защита от слишком длинных параметров
	if len(dateStr) > maxDateLen {
		writeError(w, "Слишком длинный параметр date", http.StatusBadRequest)
		return
	}
	if len(repeat) > maxRepeatLen {
		writeError(w, "Слишком длинный параметр repeat", http.StatusBadRequest)
		return
	}

	// Валидация формата date
	if len(dateStr) != 8 {
		writeError(w, "Параметр date должен быть в формате YYYYMMDD (8 символов)", http.StatusBadRequest)
		return
	}

	// Определение текущей даты (now)
	now, err := parseNow(nowStr)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызов функции NextDate для вычисления следующей даты
	next, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		log.Printf("Ошибка при вычислении следующей даты: %v", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращение следующей даты
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Отправка ответа с проверкой ошибок
	n, err := w.Write([]byte(next))
	if err != nil {
		log.Printf("Ошибка при отправке ответа: %v (отправлено байт: %d)", err, n)
	} else if n == 0 {
		log.Printf("Предупреждение: отправлено 0 байт для даты %s", next)
	}
}

// Функция parseNow парсит параметр now или возвращает текущую дату
func parseNow(nowStr string) (time.Time, error) {
	if nowStr == "" {
		// Если now не указан, используем текущую дату
		now := time.Now()
		// Обрезаем время до начала дня
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}

	// Проверка длины строки now
	if len(nowStr) != 8 {
		return time.Time{}, errors.New("параметр now должен быть в формате YYYYMMDD (8 символов)")
	}

	// Парсим переданную дату now
	now, err := time.Parse(nextdate.DateFormat, nowStr)
	if err != nil {
		return time.Time{}, err
	}

	return now, nil
}
