ALTER TABLE
    expense_tracker.accounts DROP COLUMN IF EXISTS version;

ALTER TABLE
    expense_tracker.users DROP COLUMN IF EXISTS version;

ALTER TABLE
    expense_tracker.categories DROP COLUMN IF EXISTS version;

ALTER TABLE
    expense_tracker.expenses DROP COLUMN IF EXISTS version;