package main

//Импортируем библиотеки
import (
	"fmt"
	"log"
	"net/http"
)

// Главная страница
func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Добро пожаловать на мой проект GO")
}

// Страница О проекте
func AboutPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Это проект для изучения веб-разработки на Go.")
}

// Проверка ping
func PingPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

// Активация маршрутов программы
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", MainPage)
	mux.HandleFunc("GET /about", AboutPage)
	mux.HandleFunc("GET /ping", PingPage)

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
