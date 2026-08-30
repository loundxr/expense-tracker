ALTER TABLE
    expense_tracker.users
ADD
    COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin'));