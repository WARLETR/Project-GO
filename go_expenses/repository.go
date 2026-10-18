package main

import "database/sql"

// Создаем структуру Expense
type Expense struct {
	Id          int    // айди для базы
	Amount      int    // Сумма
	Description string // Описание
	Date        string // Дата
}

// Структура репозитория как в задании
type ExpenseRepository struct {
	db *sql.DB
}

// функция добавления траты в таблицу
func (repo *ExpenseRepository) Create(e Expense) {
	// вставляем данные через знаки вопроса чтобы не было инъекций
	repo.db.Exec("INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)", e.Amount, e.Description, e.Date)
}

// Получаем все траты из бд для главной страницы
func (repo *ExpenseRepository) GetAll() []Expense {
	rows, _ := repo.db.Query("SELECT id, amount, description, date FROM expenses")

	var mas []Expense
	// перебираем строчки которые вернула база
	for rows.Next() {
		var e Expense
		rows.Scan(&e.Id, &e.Amount, &e.Description, &e.Date)
		mas = append(mas, e)
	}
	return mas
}

// Находим одну запись по ее id (нужно когда редактируем)
func (repo *ExpenseRepository) GetByID(id int) Expense {
	var e Expense
	rows, _ := repo.db.Query("SELECT id, amount, description, date FROM expenses WHERE id = ?", id)
	for rows.Next() {
		rows.Scan(&e.Id, &e.Amount, &e.Description, &e.Date)
	}
	return e
}

//Обновление записи когда поменяли в форме
func (repo *ExpenseRepository) Update(e Expense) {
	repo.db.Exec("UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?", e.Amount, e.Description, e.Date, e.Id)
}

//Удаление записи по айди
func (repo *ExpenseRepository) Delete(id int) {
	repo.db.Exec("DELETE FROM expenses WHERE id = ?", id)
}
