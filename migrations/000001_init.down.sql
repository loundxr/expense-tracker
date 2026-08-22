DROP TABLE IF EXISTS expense_tracker.expenses;

DROP TABLE IF EXISTS expense_tracker.categories;

DROP TABLE IF EXISTS expense_tracker.account_users;

DROP TABLE IF EXISTS expense_tracker.accounts;

DROP TABLE IF EXISTS expense_tracker.users;

DROP INDEX IF EXISTS expense_tracker.idx_expenses_account_id;

DROP INDEX IF EXISTS expense_tracker.idx_account_users_user_id;

DROP INDEX IF EXISTS expense_tracker.idx_date_expenses;

DROP SCHEMA IF EXISTS expense_tracker;