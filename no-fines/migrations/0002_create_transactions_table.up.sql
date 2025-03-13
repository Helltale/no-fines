CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    reserve_id INT NOT NULL REFERENCES reserves(id),
    amount DECIMAL(20, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);