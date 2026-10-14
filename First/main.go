package main

import (
	"fmt"
	"log"
	"net/http"
)
// homeHandler — главная страница с приветствием
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Страница не найдена", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Добро пожаловать на главную страницу проекта «Учёт личных трат».")
}

// aboutHandler — описание проекта
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
			http.Error(w, "Страница не найдена", http.StatusNotFound)
			return
	}
	fmt.Fprintln(w, "Это учебный проект: простой HTTP-сервер на Go с использованием net/http.")
}

// pingHandler
func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
				w.Header().Set("Allow", "GET")
				http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
				return
		}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}  

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Страница не найдена", http.StatusNotFound)
		})

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/ping", pingHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
