// pkg/db/task.go
package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

// Task представляет задачу в базе данных
type Task struct {
	ID      string `json:"id"`      // ID как строка для JSON
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Title   string `json:"title"`   // заголовок задачи (обязательное поле)
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения
}

// GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {
	// Преобразуем строковый ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("некорректный идентификатор: %s", id)
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	var task Task
	var dbID int64

	err = DB.QueryRow(query, idInt).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача с ID %s не найдена", id)
		}
		return nil, err
	}

	// Преобразуем ID в строку для JSON
	task.ID = strconv.FormatInt(dbID, 10)

	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	// Преобразуем строковый ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор: %s", task.ID)
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idInt)
	if err != nil {
		return err
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", task.ID)
	}

	return nil
}

// Tasks возвращает список задач с ограничением по количеству
// Если search параметр указан, выполняет поиск по заголовку, комментарию или дате
func Tasks(limit int, search string) ([]*Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	// Проверяем, есть ли параметр поиска
	if search == "" {
		// Простой запрос без поиска
		query := `SELECT id, date, title, comment, repeat FROM scheduler 
				  ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	} else {
		// Поиск по заголовку, комментарию или дате
		// Проверяем, является ли search датой в формате DD.MM.YYYY
		if isDateSearch(search) {
			// Преобразуем дату из формата DD.MM.YYYY в YYYYMMDD
			date := convertDateFormat(search)
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
					  WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, date, limit)
		} else {
			// Поиск по подстроке в title или comment
			searchPattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
					  WHERE title LIKE ? OR comment LIKE ? 
					  ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, searchPattern, searchPattern, limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем слайс задач (не nil, чтобы в JSON был [], а не null)
	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}
		var id int64 // временная переменная для сканирования
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		// Преобразуем int64 в строку для JSON
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddTask добавляет новую задачу в базу данных
// Возвращает ID добавленной задачи
func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (?, ?, ?, ?)
	`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteTask удаляет задачу по её ID
func DeleteTask(id string) error {
	// Преобразуем строковый ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор: %s", id)
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, idInt)
	if err != nil {
		return err
	}

	// Проверяем, была ли удалена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}

// isDateSearch проверяет, является ли строка датой в формате DD.MM.YYYY
func isDateSearch(s string) bool {
	// Проверяем длину и формат
	if len(s) != 10 {
		return false
	}
	if s[2] != '.' || s[5] != '.' {
		return false
	}
	// Проверяем, что все символы до и после точек - цифры
	for i := 0; i < 10; i++ {
		if i == 2 || i == 5 {
			continue
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// convertDateFormat преобразует дату из формата DD.MM.YYYY в YYYYMMDD
func convertDateFormat(dateStr string) string {
	// dateStr имеет формат DD.MM.YYYY
	day := dateStr[0:2]
	month := dateStr[3:5]
	year := dateStr[6:10]
	return year + month + day
}

// UpdateDate обновляет только дату задачи (для отметки о выполнении)
func UpdateDate(id string, newDate string) error {
	// Преобразуем строковый ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор: %s", id)
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, newDate, idInt)
	if err != nil {
		return err
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}
