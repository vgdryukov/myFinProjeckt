// pkg/nextdate/calc_nextDate.go
package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	maxDayInterval = 400
	maxIterations  = 1000
	maxYearSearch  = 100
	DateFormat     = "20060102"
)

// AfterNow проверяет, что дата date больше даты now
// Сравниваются только даты, без времени
func AfterNow(date, now time.Time) bool {
	// Получаем компоненты даты (год, месяц, день)
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	// Сравниваем по году
	if y1 > y2 {
		return true
	}
	if y1 < y2 {
		return false
	}

	// Годы равны - сравниваем по месяцу
	if m1 > m2 {
		return true
	}
	if m1 < m2 {
		return false
	}

	// Годы и месяцы равны - сравниваем по дню
	return d1 > d2
}

// NextDate вычисляет следующую дату для задачи по правилу повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка: правило повторять задачу не пустое
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Парсим начальную дату
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата %s: %v", dstart, err)
	}

	// Разбиваем правило повторения задачи на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("некорректный формат правила")
	}

	rule := parts[0]

	switch rule {
	case "d":
		return handleDayRule(now, date, parts)
	case "y":
		return handleYearRule(now, date)
	case "w":
		return handleWeekRule(now, date, parts)
	case "m":
		return handleMonthRule(now, date, parts)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", repeat)
	}
}

// handleDayRule обрабатывает правило "d N"
func handleDayRule(now time.Time, date time.Time, parts []string) (string, error) {
	// Проверяем, что указан интервал
	if len(parts) < 2 {
		return "", errors.New("для правила d не указан интервал дней")
	}

	// Парсим интервал
	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("интервал дней должен быть числом")
	}

	// Проверка допустимости интервала
	if interval <= 0 {
		return "", errors.New("интервал дней должен быть положительным числом")
	}
	if interval > maxDayInterval {
		return "", errors.New("интервал дней не может превышать " + strconv.Itoa(maxDayInterval))
	}

	// Поиск следующей даты
	for {
		date = date.AddDate(0, 0, interval)
		if AfterNow(date, now) {
			break
		}
		// Защита от бесконечного цикла
		if date.Year() > now.Year()+maxYearSearch {
			return "", errors.New("не удалось найти следующую дату")
		}
	}

	return date.Format(DateFormat), nil
}

// handleYearRule обрабатывает правило "y"
func handleYearRule(now time.Time, date time.Time) (string, error) {
	// Поиск следующей даты
	for {
		date = date.AddDate(1, 0, 0)
		if AfterNow(date, now) {
			break
		}
		// Защита от бесконечного цикла
		if date.Year() > now.Year()+maxYearSearch {
			return "", errors.New("не удалось найти следующую дату")
		}
	}

	return date.Format(DateFormat), nil
}

// handleWeekRule обрабатывает правило "w дни_недели"
// Примеры: "w 7" - воскресенье, "w 1,4,5" - понедельник, четверг, пятница
func handleWeekRule(now time.Time, date time.Time, parts []string) (string, error) {
	// Проверка, что указаны дни недели
	if len(parts) < 2 {
		return "", errors.New("для правила w не указаны дни недели")
	}

	// Парсим дни недели
	weekDaysStr := strings.Split(parts[1], ",")
	var weekDays []int
	for _, dayStr := range weekDaysStr {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", errors.New("дни недели должны быть числами")
		}
		// Проверка допустимости порядковых номеров дней недели (1-7, где 1=понедельник, 7=воскресенье)
		if day < 1 || day > 7 {
			return "", fmt.Errorf("недопустимый день недели: %d (допустимо 1-7)", day)
		}
		weekDays = append(weekDays, day)
	}

	// Начинаем поиск со следующего дня после date
	currentDate := date.AddDate(0, 0, 1)

	// Ограничение: поиск максимальным периодом maxIterations
	for i := 0; i < maxIterations; i++ {
		// Получаем день недели в Go (0-6)
		goWeekDay := int(currentDate.Weekday())
		// Преобразуем в наш формат (1-7)
		var ourWeekDay int
		if goWeekDay == 0 { // воскресенье
			ourWeekDay = 7
		} else {
			ourWeekDay = goWeekDay
		}

		// Проверка: есть ли этот день в списке разрешенных
		for _, allowedDay := range weekDays {
			if ourWeekDay == allowedDay {
				// Нашли подходящий день
				if AfterNow(currentDate, now) {
					return currentDate.Format(DateFormat), nil
				}
				// Если дата меньше now, продолжаем поиск
				break
			}
		}

		// Переход к следующему дню
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return "", errors.New("не удалось найти подходящую дату для правила w")
}

// handleMonthRule обрабатывает правило "m дни_месяца [месяцы]"
// Примеры:
// "m 4" - 4-е число каждого месяца
// "m 1,15,25" - 1, 15 и 25 числа каждого месяца
// "m -1" - последний день месяца
// "m -2" - предпоследний день месяца
// "m 3 1,3,6" - 3 января, марта и июня
// "m 1,-1 2,8" - первый и последний дни февраля и августа
func handleMonthRule(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("для правила m не указаны дни месяца")
	}

	// Парсим дни месяца
	daysStr := strings.Split(parts[1], ",")
	var monthDays []int
	for _, dayStr := range daysStr {
		dayStr = strings.TrimSpace(dayStr)
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", errors.New("дни месяца должны быть числами")
		}
		// ИСПРАВЛЕНО: разрешаем -2, -1 и 1-31
		if day < -2 || day > 31 || day == 0 {
			return "", fmt.Errorf("недопустимый день месяца: %d", day)
		}
		monthDays = append(monthDays, day)
	}

	// Парсим месяцы
	var months []int
	if len(parts) >= 3 {
		monthsStr := strings.Split(parts[2], ",")
		for _, monthStr := range monthsStr {
			monthStr = strings.TrimSpace(monthStr)
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				return "", errors.New("месяцы должны быть числами")
			}
			if month < 1 || month > 12 {
				return "", fmt.Errorf("недопустимый месяц: %d", month)
			}
			months = append(months, month)
		}
	}

	// Начинаем со следующего дня после date
	currentDate := date.AddDate(0, 0, 1)

	for i := 0; i < maxIterations; i++ {
		// Проверяем, подходит ли текущий месяц (если указаны конкретные месяцы)
		if len(months) > 0 {
			currentMonth := int(currentDate.Month())
			monthOk := false
			for _, allowedMonth := range months {
				if currentMonth == allowedMonth {
					monthOk = true
					break
				}
			}
			if !monthOk {
				currentDate = currentDate.AddDate(0, 0, 1)
				continue
			}
		}

		// Получаем последний день текущего месяца
		lastDay := getLastDayOfMonth(currentDate)
		currentDay := currentDate.Day()

		// Проверяем каждый разрешенный день
		for _, allowedDay := range monthDays {
			var targetDay int
			switch allowedDay {
			case -1:
				targetDay = lastDay
			case -2:
				targetDay = lastDay - 1
			default:
				targetDay = allowedDay
			}

			// ИСПРАВЛЕНИЕ: Пропускаем несуществующие дни
			if targetDay > lastDay {
				continue
			}

			if currentDay == targetDay && AfterNow(currentDate, now) {
				return currentDate.Format(DateFormat), nil
			}
		}

		// Переходим к следующему дню
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return "", errors.New("не удалось найти подходящую дату для правила m")
}

// getLastDayOfMonth возвращает последний день месяца для указанной даты
func getLastDayOfMonth(t time.Time) int {
	// Переход к первому дню следующего месяца за вычетом одного дня
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}
