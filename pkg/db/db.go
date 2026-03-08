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

// Функция Init инициализирует подключение к базе данных
// Принимает путь к файлу базы данных
// Возвращает ошибку, если что-то пошло не так
func Init(dbFile string) error {

	// Проверка существования файла базы данных
	_, err := os.Stat(dbFile)
	needCreate := os.IsNotExist(err)

	// Открываем базу данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверка подключения к базе данных
	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}

	// Если файл базы данных не существовал, создаем её таблицу и индекс
	if needCreate {
		log.Printf("База данных %s не найдена. Создаю новую...", dbFile)

		// Выполняем команды SQL для создания таблицы и индекса базы данных
		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return err
		}

		log.Println("Таблица scheduler и индекс успешно созданы")

	} else {

		// Повторная проверка: база данных существует (успешно создана)
		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return err
		}
	}

	// Сохранение подключения в глобальной переменной
	DB = db

	// Проверка: база данных существует и имеет правильную структуру
	if err := migrateDB(); err != nil {
		log.Printf("Предупреждение при проверке схемы БД: %v", err)
	}

	log.Printf("База данных %s успешно подключена", dbFile)

	return nil
}

// Функция Close закрывает соединение с базой данных
func Close() error {

	if DB != nil {
		return DB.Close()
	}

	return nil
}

// Функция GetDB возвращает текущее подключение к БД
func GetDB() *sql.DB {

	return DB
}

// Функция migrateDB проверяет и обновляет схему БД при необходимости
func migrateDB() error {

	_, err := DB.Exec(schema)

	return err
}
