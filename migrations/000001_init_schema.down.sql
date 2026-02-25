DROP INDEX IF EXISTS idx_loans_status;
DROP INDEX IF EXISTS idx_budgets_month_year;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_category;
DROP INDEX IF EXISTS idx_transactions_date;

DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS fixed_charges;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS transactions;
