package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"myfinproject/pkg/db"
	"myfinproject/pkg/server"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	WebDir string
	Port   int
	DBFile string
}

var config Config

func init() {
	// Настройка логирования в файл
	if err := setupLogging(); err != nil {
		log.Printf("Предупреждение: не удалось настроить логирование в файл: %v", err)
	}

	// Задаём значения полей структуры конфигурации сервера по умолчанию
	config = Config{
		WebDir: "./web",
		Port:   7540,
		DBFile: "scheduler.db",
	}

	// Переопределение порта из переменных окружения
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			config.Port = p
			log.Printf("Установлен порт из окружения: %d", p)
		} else {
			log.Printf("Предупреждение: некорректное значение TODO_PORT='%s', используется порт %d",
				envPort, config.Port)
		}
	}

	// Переопределение пути к файлу базы данных
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		config.DBFile = envDBFile
		log.Printf("Установлен путь к БД из окружения: %s", config.DBFile)
	}

	// Переопределение webDir из окружения
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

// Функция setupLogging настраивает запись логов в файл
func setupLogging() error {
	// Определяем путь к логам (сначала проверяем переменную окружения)
	logPath := os.Getenv("LOG_FILE")
	if logPath == "" {
		logPath = "/app/logs/app.log" // абсолютный путь внутри контейнера
	}

	// Создаем директорию для логов, если её нет
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию логов %s: %w", logDir, err)
	}

	// Открываем файл для логов
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл логов %s: %w", logPath, err)
	}

	// Настраиваем вывод логов одновременно в консоль и в файл
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Логирование настроено в файл: %s", logPath)
	return nil
}

func main() {
	log.Printf("Запуск с конфигурацией: Порт=%d, БД=%s, WebDir=%s",
		config.Port, config.DBFile, config.WebDir)

	if err := db.Init(config.DBFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}()

	srv := server.NewServer(config.WebDir, config.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil {
			errCh <- err
		}
	}()

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("\n✅ Сервер успешно запущен!\n")
	fmt.Printf("📁 Веб-директория: %s\n", config.WebDir)
	fmt.Printf("🔌 Порт: %d\n", config.Port)
	fmt.Printf("💾 База данных: %s\n", config.DBFile)
	fmt.Printf("📝 Логи сохраняются в: /app/logs/app.log (в контейнере)\n")

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		fmt.Println("\n🛑 Получен сигнал завершения. Останавливаем сервер...")
	case err := <-errCh:
		log.Printf("Ошибка сервера: %v", err)
		fmt.Println("\n🛑 Сервер остановлен из-за ошибки")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	fmt.Println("✅ Сервер успешно остановлен")
}
