// pkg/server/server.go
// Управление HTTP-сервером
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"myfinproject/pkg/api"
)

// Структура конфигурации сервера
type Server struct {
	webDir string
	port   int
	http   *http.Server
}

// Функция NewServer создает новый экземпляр сервера
func NewServer(webDir string, port int) *Server {
	return &Server{
		webDir: webDir,
		port:   port,
	}
}

// Функция Start запускает сервер
func (s *Server) Start() error {
	// Проверка существования папки web
	if _, err := os.Stat(s.webDir); os.IsNotExist(err) {
		log.Printf("Предупреждение: папка %s не найдена", s.webDir)
	}

	// Инициализация api до запуска сервера
	api.Init()

	// Создание файлового сервера для статических файлов
	fileServer := http.FileServer(http.Dir(s.webDir))
	http.Handle("/", fileServer)

	addr := fmt.Sprintf(":%d", s.port)

	// Создание HTTP сервера
	s.http = &http.Server{
		Addr:         addr,
		Handler:      nil, // используем DefaultServeMux
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
	fmt.Printf("Обслуживаются файлы из директории: %s\n", s.webDir)
	fmt.Printf("API доступно по адресу: http://localhost:%d/api/nextdate\n", s.port)

	// Запуск сервера
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("ошибка при запуске сервера: %v", err)
	}

	return nil
}

// Функция Stop останавливает сервер gracefully
func (s *Server) Stop(ctx context.Context) error {
	if s.http == nil {
		return nil
	}

	log.Println("Останавливаем сервер...")
	return s.http.Shutdown(ctx)
}

// Функция GetPort возвращает текущий порт сервера
func (s *Server) GetPort() int {
	return s.port
}
