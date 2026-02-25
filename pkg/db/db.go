// pkg/db/db.go
package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" // драйвер SQLite
)

// Глобальная переменная для доступа к БД
var DB *sql.DB

// Schema содержит SQL команды для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Init инициализирует подключение к базе данных
// Принимает путь к файлу базы данных
// Возвращает ошибку, если что-то пошло не так
func Init(dbFile string) error {
	// Проверяем существование файла базы данных
	_, err := os.Stat(dbFile)
	needCreate := os.IsNotExist(err)

	// Открываем базу данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}

	// Если файл не существовал, создаем таблицу и индекс
	if needCreate {
		log.Printf("База данных %s не найдена. Создаю новую...", dbFile)

		// Выполняем SQL для создания таблицы и индекса
		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return err
		}

		log.Println("Таблица scheduler и индекс успешно созданы")
	} else {
		// Проверяем, существует ли таблица (на всякий случай)
		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return err
		}
	}

	// Сохраняем подключение в глобальной переменной
	DB = db

	log.Printf("База данных %s успешно подключена", dbFile)
	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB возвращает текущее подключение к БД
func GetDB() *sql.DB {
	return DB
}
