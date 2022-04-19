CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    order_id INTEGER,
    sum BIGINT,
    processed_at TIMESTAMP DEFAULT NOW()
);