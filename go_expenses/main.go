package main

// Импортируем библиотеки
import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

// создаем глобальную переменную репозитория чтоб была видна везде
var repo ExpenseRepository

// Главная страница
func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//Если пользователь отправил форму
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
			return
		}

		//Достаем значения из полей HTML формы
		amountStr := r.FormValue("amount")
		description := r.FormValue("description")
		date := r.FormValue("date")

		//Переводим строковую сумму в число
		amount, _ := strconv.Atoi(amountStr)

		//Создаем новый объект траты
		newExpense := Expense{
			Amount:      amount,
			Description: description,
			Date:        date,
		}

		//Записываем трату в базу через репозиторий
		repo.Create(newExpense)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Вытаскиваем список трат из базы
	expensesList := repo.GetAll()

	//подключаем шаблон через пакет html/template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов", http.StatusInternalServerError)
		return
	}

	//Передаем наш массив трат внутрь HTML шаблона
	tmpl.ExecuteTemplate(w, "layout", expensesList)
}

// Функция для удаления расхода при нажатии на ссылку
func DeletePage(w http.ResponseWriter, r *http.Request) {
	//берем id из адресной строки браузера
	id_str := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(id_str)

	repo.Delete(id)

	//возвращаем обратно на главную
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Страница для редактирования расхода
func EditPage(w http.ResponseWriter, r *http.Request) {
	//если нажали кнопку сохранить в форме редактирования
	if r.Method == "POST" {
		r.ParseForm()

		id, _ := strconv.Atoi(r.FormValue("id"))
		amount, _ := strconv.Atoi(r.FormValue("amount"))
		description := r.FormValue("description")
		date := r.FormValue("date")

		exp := Expense{
			Id:          id,
			Amount:      amount,
			Description: description,
			Date:        date,
		}

		repo.Update(exp)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//если просто перешли на страницу редактирования, ищем запись по id
	id_str := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(id_str)

	exp := repo.GetByID(id)

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/edit.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", exp)
}

// Страница О проекте
func AboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/about.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов", http.StatusInternalServerError)
		return
	}
	//Данные для этой страницы не нужны, передаем nil
	tmpl.ExecuteTemplate(w, "layout", nil)
}

// Проверка ping
func PingPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// Активация маршрутов программы
func main() {
	// подключаемся к файлу базы данных
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		log.Fatal(err)
	}

	//создаем таблицу если она еще не создана
	db.Exec("CREATE TABLE IF NOT EXISTS expenses (id INTEGER PRIMARY KEY AUTOINCREMENT, amount INTEGER, description TEXT, date TEXT)")

	repo = ExpenseRepository{db: db}

	mux := http.NewServeMux()

	mux.HandleFunc("/", MainPage)
	mux.HandleFunc("/delete", DeletePage)
	mux.HandleFunc("/edit", EditPage)
	mux.HandleFunc("GET /about", AboutPage)
	mux.HandleFunc("GET /ping", PingPage)

	log.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
