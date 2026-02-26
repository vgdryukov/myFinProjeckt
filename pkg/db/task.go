// pkg/db/task.go
package db

// Task представляет задачу в базе данных
type Task struct {
	ID      int64  `json:"id"`      // идентификатор (int64 для БД)
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Title   string `json:"title"`   // заголовок задачи (обязательное поле)
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения
}
