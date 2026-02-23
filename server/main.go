package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	webDir string = "../web" // Относительный путь от server/ к web/
	port   int    = 7540     // Порт для прослушивания
)

func init() {
	// Проверяем переменную окружения TODO_PORT
	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Проверяем, существует ли папка web
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Printf("Предупреждение: папка %s не найдена", webDir)
	}

	// Запускаем сервер при инициализации пакета
	go func() {
		// Создаем файловый сервер для директории web
		fileServer := http.FileServer(http.Dir(webDir))

		// Обрабатываем все запросы с помощью файлового сервера
		http.Handle("/", fileServer)

		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
		fmt.Printf("Обслуживаются файлы из директории: %s\n", webDir)

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("Ошибка при запуске сервера: %v", err)
		}
	}()

	// Даем серверу время на запуск
	time.Sleep(1 * time.Second)
}

func main() {
	// В обычном режиме работаем бесконечно
	fmt.Println("Сервер работает. Нажмите Ctrl+C для остановки")
	select {}
}
