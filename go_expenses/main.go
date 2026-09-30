package main

// Импортируем библиотеки
import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Создаем структуру Expense
type Expense struct {
	Amount      int    // Сумма
	Description string // Описание
	Date        string // Дата
}

// Временное хранилище в памяти
var expensesList = []Expense{}

// Главная страница
func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Если пользователь отправил форму
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
			return
		}

		// Достаем значения из полей HTML-формы
		amountStr := r.FormValue("amount")
		description := r.FormValue("description")
		date := r.FormValue("date")

		// Переводим строковую сумму в число
		amount, _ := strconv.Atoi(amountStr)

		// Создаем новый объект траты
		newExpense := Expense{
			Amount:      amount,
			Description: description,
			Date:        date,
		}

		// Добавляем созданную трату в наш массив в памяти
		expensesList = append(expensesList, newExpense)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Компилируем файлы шаблонов через пакет html/template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов", http.StatusInternalServerError)
		return
	}

	// Передаем наш массив трат внутрь HTML-шаблона
	tmpl.ExecuteTemplate(w, "layout", expensesList)
}

// Страница О проекте
func AboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/about.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов", http.StatusInternalServerError)
		return
	}
	// Данные для этой страницы не нужны, передаем nil
	tmpl.ExecuteTemplate(w, "layout", nil)
}

// Проверка ping
func PingPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// Активация маршрутов программы
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", MainPage)
	mux.HandleFunc("GET /about", AboutPage)
	mux.HandleFunc("GET /ping", PingPage)

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
