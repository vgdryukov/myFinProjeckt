// pkg/db/task_db.go
// Операции с задачами в базе данных
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

// Функция GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {

	// Преобразование строкового ID в int64 для запроса к БД
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

	// Преобразование ID int64 в строку для JSON
	task.ID = strconv.FormatInt(dbID, 10)

	return &task, nil
}

// Функция UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {

	// Преобразование строкового ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор: %s", task.ID)
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idInt)
	if err != nil {
		return err
	}

	// Проверка: была ли обновлена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", task.ID)
	}

	return nil
}

// Функция Tasks возвращает список задач с ограничением по количеству
// Если search параметр указан, выполняет поиск по заголовку, комментарию или дате
func Tasks(limit int, search string) ([]*Task, error) {

	var (
		rows *sql.Rows
		err  error
	)

	// Проверка: есть ли параметр поиска search
	if search == "" {
		// Простой запрос без поиска
		query := `SELECT id, date, title, comment, repeat FROM scheduler 
				  ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)

	} else {

		// Поиск по заголовку, комментарию или дате
		// Проверка: является ли search датой в формате DD.MM.YYYY
		if isDateSearch(search) {

			// Преобразование даты из формата DD.MM.YYYY в YYYYMMDD
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

	// Создание слайса задач (не nil, чтобы в JSON был [], а не null)
	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}
		var id int64 // временная переменная для сканирования
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		// Преобразование ID int64 в строку для JSON
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Функция AddTask добавляет новую задачу в базу данных
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

// Функция DeleteTask удаляет задачу по её ID
func DeleteTask(id string) error {

	// Преобразование строкового ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор: %s", id)
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, idInt)
	if err != nil {
		return err
	}

	// Проверка: была ли удалена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}

// Функция isDateSearch проверяет, является ли строка датой в формате DD.MM.YYYY
func isDateSearch(s string) bool {

	// Проверка длины и формата строки даты
	if len(s) != 10 {
		return false
	}
	if s[2] != '.' || s[5] != '.' {
		return false
	}

	// Проверка: все символы до и после точек - цифры
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

// Функция convertDateFormat преобразует дату из формата DD.MM.YYYY в YYYYMMDD
func convertDateFormat(dateStr string) string {

	// dateStr имеет формат DD.MM.YYYY
	day := dateStr[0:2]
	month := dateStr[3:5]
	year := dateStr[6:10]

	return year + month + day
}

// Функция UpdateDate обновляет только дату задачи (для отметки о выполнении)
func UpdateDate(id string, newDate string) error {

	// Преобразование строкового ID в int64 для запроса к БД
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор: %s", id)
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, newDate, idInt)
	if err != nil {
		return err
	}

	// Проверка: была ли обновлена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}
