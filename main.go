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
		}
	}

	// Переопределяем базу данных из переменных окружения
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		config.DBFile = envDBFile
		log.Printf("Используется путь к БД из переменной окружения: %s", config.DBFile)
	}
}

func main() {
	// Инициализируем базу данных
	if err := db.Init(config.DBFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Создаем и запускаем сервер
	srv := server.NewServer(config.WebDir, config.Port)
	srv.Start()

	fmt.Printf("Конфигурация:\n")
	fmt.Printf("  Порт: %d\n", config.Port)
	fmt.Printf("  База данных: %s\n", config.DBFile)
	fmt.Printf("  Веб-директория: %s\n", config.WebDir)
	fmt.Println("Нажмите Ctrl+C для остановки")

	// Грациозное завершение
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nЗавершение работы...")
}
