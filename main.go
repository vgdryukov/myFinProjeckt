// main.go
package main

import (
	"fmt"
	"log"
	"myfinproject/pkg/db"
	"myfinproject/pkg/server"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Config содержит все настройки приложения
type Config struct {
	WebDir string
	Port   int
	DBFile string
}

var config Config

func init() {
	// Значения по умолчанию
	config = Config{
		WebDir: "./web",
		Port:   7540,
		DBFile: "scheduler.db",
	}

	// Переопределяем порт из переменных окружения
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			config.Port = p
			log.Printf("Установлен порт из окружения: %d", p)
		} else {
			log.Printf("Предупреждение: некорректное значение TODO_PORT='%s', используется порт %d",
				envPort, config.Port)
		}
	}

	// Переопределяем базу данных из переменных окружения
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		config.DBFile = envDBFile
		log.Printf("Установлен путь к БД из окружения: %s", config.DBFile)
	}

	// Можно также переопределить webDir из окружения (опционально)
	if envWebDir := os.Getenv("TODO_WEB_DIR"); envWebDir != "" {
		config.WebDir = envWebDir
		log.Printf("Установлена веб-директория из окружения: %s", config.WebDir)
	}
}

func main() {
	// Выводим информацию о конфигурации
	log.Printf("Запуск с конфигурацией: Порт=%d, БД=%s, WebDir=%s",
		config.Port, config.DBFile, config.WebDir)

	// Инициализируем базу данных
	if err := db.Init(config.DBFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}()

	// Создаем и запускаем сервер
	srv := server.NewServer(config.WebDir, config.Port)
	srv.Start()

	fmt.Printf("\nСервер успешно запущен!\n")
	fmt.Printf("Веб-директория: %s\n", config.WebDir)
	fmt.Printf("Порт: %d\n", config.Port)
	fmt.Printf("База данных: %s\n", config.DBFile)
	fmt.Printf("API доступно:\n")
	fmt.Printf(" - GET  /api/nextdate?now=&date=&repeat=\n")
	fmt.Printf(" - POST /api/task\n")
	fmt.Printf(" - GET  /api/task?id= \n")
	fmt.Printf(" - PUT  /api/task\n")
	fmt.Printf(" - GET  /api/tasks?search= \n")

	// Грациозное завершение
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nПолучен сигнал завершения. Останавливаем сервер...")
	fmt.Println("Сервер успешно остановлен")
}
