package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

var repo ExpenseRepository

// Главная страница: обработка формы POST и вывод списка с фильтрами GET
func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Обработка отправки формы добавления траты
	if r.Method == http.MethodPost {
		r.ParseForm()
		amount, _ := strconv.Atoi(r.FormValue("amount"))
		description := r.FormValue("description")
		date := r.FormValue("date")
		catID, _ := strconv.Atoi(r.FormValue("category_id"))

		newExpense := Expense{
			Amount:      amount,
			Description: description,
			Date:        date,
			CategoryId:  catID,
		}

		repo.Create(newExpense)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Считываем параметры фильтрации из адресной строки
	catFilter := r.URL.Query().Get("category")
	fromFilter := r.URL.Query().Get("date_from")
	toFilter := r.URL.Query().Get("date_to")

	expensesList := repo.GetAllFiltered(catFilter, fromFilter, toFilter)
	categoriesList := repo.GetCategories()

	data := struct {
		Expenses    []Expense
		Categories  []Category
		SelectedCat string
		DateFrom    string
		DateTo      string
	}{
		Expenses:    expensesList,
		Categories:  categoriesList,
		SelectedCat: catFilter,
		DateFrom:    fromFilter,
		DateTo:      toFilter,
	}

	// Генерация интерфейса главной страницы через шаблонизатор
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", 500)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

// Страница удаления записи по id
func DeletePage(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(id_str)
	repo.Delete(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Страница редактирования траты (GET - показать форму со старыми данными, POST - сохранить)
func EditPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		amount, _ := strconv.Atoi(r.FormValue("amount"))
		description := r.FormValue("description")
		date := r.FormValue("date")
		catID, _ := strconv.Atoi(r.FormValue("category_id"))

		exp := Expense{
			Id:          id,
			Amount:      amount,
			Description: description,
			Date:        date,
			CategoryId:  catID,
		}

		repo.Update(exp)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id_str := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(id_str)
	exp := repo.GetByID(id)
	categoriesList := repo.GetCategories()

	data := struct {
		Exp        Expense
		Categories []Category
	}{
		Exp:        exp,
		Categories: categoriesList,
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/edit.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", 500)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

// Страница "О проекте"
func AboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/about.html")
	if err != nil {
		http.Error(w, "Ошибка", 500)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}

// Страница Ping-Pong проверки
func PingPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func main() {
	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		log.Fatal(err)
	}

	// Автоматическое создание таблиц, если их нет
	db.Exec("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	db.Exec("CREATE TABLE IF NOT EXISTS expenses (id INTEGER PRIMARY KEY AUTOINCREMENT, amount INTEGER, description TEXT, date TEXT, category_id INTEGER)")

	// Заполнение дефолтных категорий, если таблица пустая
	var count int
	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO categories (name) VALUES ('Еда')")
		db.Exec("INSERT INTO categories (name) VALUES ('Транспорт')")
		db.Exec("INSERT INTO categories (name) VALUES ('Развлечения')")
		db.Exec("INSERT INTO categories (name) VALUES ('ЖКХ')")
	}

	repo = ExpenseRepository{db: db}

	// Создание и настройка маршрутизатора
	mux := http.NewServeMux()
	mux.HandleFunc("/", MainPage)
	mux.HandleFunc("/delete", DeletePage)
	mux.HandleFunc("/edit", EditPage)
	mux.HandleFunc("GET /about", AboutPage)
	mux.HandleFunc("GET /ping", PingPage)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
