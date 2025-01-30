CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    balance DECIMAL(20,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_accounts_created_at ON accounts(created_at);