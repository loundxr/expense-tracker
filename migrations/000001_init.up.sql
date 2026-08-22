CREATE SCHEMA IF NOT EXISTS expense_tracker;

CREATE TABLE
    IF NOT EXISTS expense_tracker.users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) NOT NULL UNIQUE CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS expense_tracker.accounts (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL CHECK (char_length(name) BETWEEN 2 AND 100),
        user_id INT NOT NULL REFERENCES expense_tracker.users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS expense_tracker.account_users (
        account_id INT NOT NULL REFERENCES expense_tracker.accounts (id) ON DELETE CASCADE,
        user_id INT NOT NULL REFERENCES expense_tracker.users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (account_id, user_id)
    );

CREATE TABLE
    IF NOT EXISTS expense_tracker.categories (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) CHECK (char_length(name) BETWEEN 2 AND 100),
        user_id INT NOT NULL REFERENCES expense_tracker.users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS expense_tracker.expenses (
        id SERIAL PRIMARY KEY,
        account_id INT NOT NULL REFERENCES expense_tracker.accounts (id) ON DELETE CASCADE,
        user_id INT NOT NULL REFERENCES expense_tracker.users (id) ON DELETE CASCADE,
        category_id INT REFERENCES expense_tracker.categories (id) ON DELETE SET NULL,
        amount BIGINT NOT NULL CHECK (amount >= 0),
        currency VARCHAR(3) NOT NULL DEFAULT 'USD' CHECK (currency ~ '^[A-Z]{3}$'),
        description TEXT CHECK (char_length(description) <= 500),
        date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_expenses_account_id ON expense_tracker.expenses (account_id);

CREATE INDEX idx_account_users_user_id ON expense_tracker.account_users (user_id);

CREATE INDEX idx_date_expenses ON expense_tracker.expenses (date);