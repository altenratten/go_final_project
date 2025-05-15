package server

import (
	"go1f/pkg/api"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	// Значение порта по умолчанию
	port := "7540"

	// Если переменная окружения TODO_PORT задана — используем её
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	api.Init()

	// Файловый сервер для директории web
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
