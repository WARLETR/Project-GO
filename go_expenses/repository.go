package main

import "database/sql"

// Структура одной траты с учетом категории
type Expense struct {
	Id           int
	Amount       int
	Description  string
	Date         string
	CategoryId   int
	CategoryName string
}

// Структура категории расходов
type Category struct {
	Id   int
	Name string
}

// Репозиторий для изоляции всех SQL-запросов к базе данных
type ExpenseRepository struct {
	db *sql.DB
}

// Достаем список всех категорий для выпадающего списка в форме
func (repo *ExpenseRepository) GetCategories() []Category {
	rows, err := repo.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var cats []Category
	for rows.Next() {
		var c Category
		rows.Scan(&c.Id, &c.Name)
		cats = append(cats, c)
	}
	return cats
}

// Добавление новой траты в базу через защищенные плейсхолдеры
func (repo *ExpenseRepository) Create(e Expense) {
	repo.db.Exec("INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)",
		e.Amount, e.Description, e.Date, e.CategoryId)
}

// Получение трат с фильтрацией по категории и датам через LEFT JOIN
func (repo *ExpenseRepository) GetAllFiltered(catID string, dateFrom string, dateTo string) []Expense {
	query := `
		SELECT e.id, e.amount, e.description, e.date, COALESCE(c.name, 'Без категории') 
		FROM expenses e 
		LEFT JOIN categories c ON e.category_id = c.id 
		WHERE 1=1
	`
	var args []interface{}

	if catID != "" && catID != "0" {
		query += " AND e.category_id = ?"
		args = append(args, catID)
	}

	if dateFrom != "" {
		query += " AND e.date >= ?"
		args = append(args, dateFrom)
	}

	if dateTo != "" {
		query += " AND e.date <= ?"
		args = append(args, dateTo)
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var mas []Expense
	for rows.Next() {
		var e Expense
		rows.Scan(&e.Id, &e.Amount, &e.Description, &e.Date, &e.CategoryName)
		mas = append(mas, e)
	}
	return mas
}

// Поиск конкретной траты по ее id для страницы редактирования
func (repo *ExpenseRepository) GetByID(id int) Expense {
	var e Expense
	rows, err := repo.db.Query("SELECT id, amount, description, date, category_id FROM expenses WHERE id = ?", id)
	if err != nil {
		return e
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&e.Id, &e.Amount, &e.Description, &e.Date, &e.CategoryId)
	}
	return e
}

// Изменение данных существующей траты в базе по id
func (repo *ExpenseRepository) Update(e Expense) {
	repo.db.Exec("UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?",
		e.Amount, e.Description, e.Date, e.CategoryId, e.Id)
}

// Полное удаление записи из базы по ее id
func (repo *ExpenseRepository) Delete(id int) {
	repo.db.Exec("DELETE FROM expenses WHERE id = ?", id)
}
