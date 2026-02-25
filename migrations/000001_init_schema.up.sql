CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'MAD',
    category_name TEXT NOT NULL,
    category_type TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS budgets (
    id TEXT PRIMARY KEY,
    category_name TEXT NOT NULL,
    limit_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'MAD',
    month INTEGER NOT NULL CHECK(month >= 1 AND month <= 12),
    year INTEGER NOT NULL,
    UNIQUE(category_name, month, year)
);

CREATE TABLE IF NOT EXISTS fixed_charges (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    amount_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'MAD',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS loans (
    id TEXT PRIMARY KEY,
    lender_name TEXT NOT NULL,
    amount_cents INTEGER NOT NULL,
    amount_paid_cents INTEGER NOT NULL DEFAULT 0,
    currency TEXT NOT NULL DEFAULT 'MAD',
    borrowed_at DATETIME NOT NULL,
    paid_back_at DATETIME,
    status TEXT NOT NULL CHECK(status IN ('active', 'paid_back')),
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_name);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_budgets_month_year ON budgets(month, year);
CREATE INDEX IF NOT EXISTS idx_loans_status ON loans(status);
