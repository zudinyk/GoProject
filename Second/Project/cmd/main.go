package main

import (
	"log"
	"net/http"
	"Project/internal/handlers"
) 

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /expenses", handlers.ExpensesHandler)
	mux.HandleFunc("GET /expenses/add", handlers.AddHandler)
	mux.HandleFunc("POST /expenses/add", handlers.AddHandler)

	// Статика (CSS)
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Редирект с главной на /expenses
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	})

	log.Println("Сервер запущен: http://localhost:8080")
	// Запуск сервера
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
