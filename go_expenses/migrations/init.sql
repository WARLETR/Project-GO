-- таблица категорий
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);

-- базовые категории
INSERT INTO categories (name) VALUES ('Еда');
INSERT INTO categories (name) VALUES ('Транспорт');
INSERT INTO categories (name) VALUES ('Развлечения');
INSERT INTO categories (name) VALUES ('ЖКХ');

-- таблица расходов с привязкой к категории
CREATE TABLE expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount INTEGER,
    description TEXT,
    date TEXT,
    category_id INTEGER
);