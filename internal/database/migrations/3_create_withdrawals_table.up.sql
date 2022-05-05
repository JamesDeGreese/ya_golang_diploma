CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    order_number VARCHAR(255),
    sum BIGINT,
    processed_at TIMESTAMP DEFAULT NOW()
);