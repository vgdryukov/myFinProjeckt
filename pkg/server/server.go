// pkg/server/server.go
package server

import (
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
}

// Функция NewServer создает новый экземпляр сервера
func NewServer(webDir string, port int) *Server {

	return &Server{
		webDir: webDir,
		port:   port,
	}
}

// Функция Start запускает сервер в фоновом режиме
func (s *Server) Start() {

	go func() {

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
		fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
		fmt.Printf("Обслуживаются файлы из директории: %s\n", s.webDir)
		fmt.Printf("API доступно по адресу: http://localhost:%d/api/nextdate\n", s.port)

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("Ошибка при запуске сервера: %v", err)
		}
	}()

	// Приостанавливаем работу на 1 секунду, чтобы дать серверу время на запуск
	time.Sleep(1 * time.Second)

}

// Функция GetPort возвращает текущий порт сервера
func (s *Server) GetPort() int {

	return s.port
}
