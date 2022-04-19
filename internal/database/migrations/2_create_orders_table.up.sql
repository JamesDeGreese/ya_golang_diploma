CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    number VARCHAR(255),
    status VARCHAR(255),
    accrual INTEGER,
    uploaded_at TIMESTAMP DEFAULT NOW()
);