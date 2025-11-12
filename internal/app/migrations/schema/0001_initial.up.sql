CREATE TABLE IF NOT EXISTS prices (
    unique_key TEXT,
    id INTEGER,
    name TEXT,
    category TEXT,
    price DECIMAL(10, 2),
    create_date DATE,
    CONSTRAINT unique_record PRIMARY KEY (unique_key, name, category, create_date)
);

CREATE TABLE IF NOT EXISTS report (
    id INTEGER,
    name TEXT,
    category TEXT,
    price DECIMAL(10, 2),
    create_date DATE,
    CONSTRAINT report_record PRIMARY KEY (name, category, create_date)
);
