ALTER TABLE
    expense_tracker.accounts
ADD
    COLUMN version INT NOT NULL DEFAULT 1;

ALTER TABLE
    expense_tracker.users
ADD
    COLUMN version INT NOT NULL DEFAULT 1;

ALTER TABLE
    expense_tracker.categories
ADD
    COLUMN version INT NOT NULL DEFAULT 1;

ALTER TABLE
    expense_tracker.expenses
ADD
    COLUMN version INT NOT NULL DEFAULT 1;