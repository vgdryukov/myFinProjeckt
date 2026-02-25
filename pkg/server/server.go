// pkg/server/server.go
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"myfinproject/pkg/api" // импортируем пакет api
)

type Server struct {
	webDir string
	port   int
}

// NewServer создает новый экземпляр сервера
func NewServer(webDir string, port int) *Server {
	return &Server{
		webDir: webDir,
		port:   port,
	}
}

// Start запускает сервер в фоновом режиме
func (s *Server) Start() {
	go func() {
		// Проверяем существование папки web
		if _, err := os.Stat(s.webDir); os.IsNotExist(err) {
			log.Printf("Предупреждение: папка %s не найдена", s.webDir)
		}

		// ИНИЦИАЛИЗИРУЕМ API ПЕРЕД ЗАПУСКОМ СЕРВЕРА
		api.Init()

		// Создаем файловый сервер для статических файлов
		fileServer := http.FileServer(http.Dir(s.webDir))
		http.Handle("/", fileServer)

		addr := fmt.Sprintf(":%d", s.port)
		fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
		fmt.Printf("Обслуживаются файлы из директории: %s\n", s.webDir)
		fmt.Printf("API доступно по адресу: http://localhost:%d/api/nextdate\n", s.port)

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("Ошибка при запуске сервера: %v", err)
		}
	}()

	// Даем серверу время на запуск
	time.Sleep(1 * time.Second)
}

// GetPort возвращает текущий порт сервера
func (s *Server) GetPort() int {
	return s.port
}
