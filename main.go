// Главный файл приложения, "точка входа"
package main

import (
	"context"
	"fmt"
	"log"
	"myfinproject/pkg/db"
	"myfinproject/pkg/server"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Структура Config содержит настройки для работы сервера
type Config struct {
	WebDir string
	Port   int
	DBFile string
}

var config Config

func init() {

	// Задаём значения полей структуры конфигурации сервера по умолчанию
	config = Config{
		WebDir: "./web",
		Port:   7540,
		DBFile: "scheduler.db",
	}

	// Переопределение порта из переменных окружения, если эта переменая определена
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			config.Port = p
			log.Printf("Установлен порт из окружения: %d", p)
		} else {
			log.Printf("Предупреждение: некорректное значение TODO_PORT='%s', используется порт %d",
				envPort, config.Port)
		}
	}

	// Переопределение путя к файлу базы данных из переменных окружения, если эта переменая определена
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		config.DBFile = envDBFile
		log.Printf("Установлен путь к БД из окружения: %s", config.DBFile)
	}

	// Переопределение webDir из окружения (опционально), если эта переменая определена
	if envWebDir := os.Getenv("TODO_WEB_DIR"); envWebDir != "" {
		config.WebDir = envWebDir
		log.Printf("Установлена веб-директория из окружения: %s", config.WebDir)
	}

	// Вывод в лог информации об аутентификации
	if envPassword := os.Getenv("TODO_PASSWORD"); envPassword != "" {
		log.Printf("Аутентификация включена (требуется пароль)")
	} else {
		log.Printf("Аутентификация отключена (пароль не установлен)")
	}

}

func main() {

	// Вывод в лог информации о конфигурации работы сервера
	log.Printf("Запуск с конфигурацией: Порт=%d, БД=%s, WebDir=%s",
		config.Port, config.DBFile, config.WebDir)

	// Инициализация базы данных
	if err := db.Init(config.DBFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}()

	// Создание и запуск сервера из пакета server
	srv := server.NewServer(config.WebDir, config.Port)

	// Запускаем сервер и получаем канал для ошибок
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil {
			errCh <- err
		}
	}()

	// Даем серверу немного времени на запуск
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("\n✅ Сервер успешно запущен!\n")
	fmt.Printf("📁 Веб-директория: %s\n", config.WebDir)
	fmt.Printf("🔌 Порт: %d\n", config.Port)
	fmt.Printf("💾 База данных: %s\n", config.DBFile)

	// Вывод в консоль информации об аутентификации
	if envPassword := os.Getenv("TODO_PASSWORD"); envPassword != "" {
		fmt.Printf("🔐 Аутентификация: ВКЛЮЧЕНА (требуется пароль)\n")
	} else {
		fmt.Printf("🔓 Аутентификация: ОТКЛЮЧЕНА\n")
	}

	fmt.Printf("\n📝 API доступно:\n")
	fmt.Printf("    - POST /api/signin (для аутентификации)\n")
	fmt.Printf("    - GET  /api/nextdate?now=&date=&repeat=\n")
	fmt.Printf("    - POST /api/task\n")
	fmt.Printf("    - GET  /api/task?id= \n")
	fmt.Printf("    - PUT  /api/task\n")
	fmt.Printf("    - DELETE /api/task?id= \n")
	fmt.Printf("    - POST /api/task/done?id= \n")
	fmt.Printf("    - GET  /api/tasks?search= \n")

	fmt.Printf("\n🛑 Нажмите Ctrl+C для остановки сервера\n")

	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем либо сигнал ОС, либо ошибку сервера
	select {
	case <-quit:
		fmt.Println("\n🛑 Получен сигнал завершения. Останавливаем сервер...")
	case err := <-errCh:
		log.Printf("Ошибка сервера: %v", err)
		fmt.Println("\n🛑 Сервер остановлен из-за ошибки")
	}

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем сервер gracefully
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	fmt.Println("✅ Сервер успешно остановлен")
}
