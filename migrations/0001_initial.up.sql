CREATE TABLE IF NOT EXISTS prices (
    id INTEGER PRIMARY KEY,
    name TEXT,
    category TEXT,
    price DECIMAL(10, 2),
    create_date DATE
);
