// pkg/api/api.go
package api

import (
	"net/http"
)

// Константа для формата даты (используется во многих местах)
const DateFormat = "20060102"

// Init регистрирует все API обработчики
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	// Здесь будут регистрироваться другие обработчики по мере добавления
}
