-- таблица для расходов
CREATE TABLE expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount INTEGER,
    description TEXT,
    date TEXT
);