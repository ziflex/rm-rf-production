CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    document_number VARCHAR(11) UNIQUE NOT NULL
);

CREATE TYPE operation_type AS ENUM ('purchase', 'installment_purchase', 'withdrawal', 'payment');

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id) NOT NULL,
    operation_type operation_type NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);