// pkg/main.go
package main

import (
	"fmt"
	"myfinproject/pkg/server"
)

var (
	webDir string = "./web"
	port   int    = 7540
	srv    *server.Server
)

func init() {
	// Создаем и запускаем сервер при инициализации
	srv = server.NewServer(webDir, port)
	srv.Start()
}

func main() {
	fmt.Println("Сервер работает. Нажмите Ctrl+C для остановки")
	select {}
}
