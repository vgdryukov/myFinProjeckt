// server.go
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	webDir string
	port   int
}

func NewServer(webDir string, port int) *Server {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	return &Server{
		webDir: webDir,
		port:   port,
	}
}

func (s *Server) Start() {
	go func() {
		if _, err := os.Stat(s.webDir); os.IsNotExist(err) {
			log.Printf("Предупреждение: папка %s не найдена", s.webDir)
		}
		fileServer := http.FileServer(http.Dir(s.webDir))
		http.Handle("/", fileServer)
		addr := fmt.Sprintf(":%d", s.port)
		fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
		fmt.Printf("Обслуживаются файлы из директории: %s\n", s.webDir)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("Ошибка при запуске сервера: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)
}

func (s *Server) GetPort() int {
	return s.port
}
