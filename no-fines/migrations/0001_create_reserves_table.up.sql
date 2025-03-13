CREATE TABLE reserves (
    id SERIAL PRIMARY KEY,
    currency VARCHAR(10) NOT NULL,
    amount DECIMAL(20, 2) NOT NULL
);

-- test
INSERT INTO reserves (currency, amount) VALUES ('RUB', 1000000.00);
INSERT INTO reserves (currency, amount) VALUES ('USD', 100000.00);
INSERT INTO reserves (currency, amount) VALUES ('EUR', 50000.00);