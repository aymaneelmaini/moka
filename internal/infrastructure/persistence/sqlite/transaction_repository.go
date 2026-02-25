package sqlite

import (
	"database/sql"
	"fmt"
	"moka/internal/domain/transaction"
	"moka/internal/shared"
	"time"
)

type TransactionRepository struct {
	db *DB
}

func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(tx transaction.Transaction) error {
	query := `
		INSERT INTO transactions (id, amount_cents, currency, category_name, category_type, description, type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		tx.ID(),
		tx.Amount().AmountInCents(),
		tx.Amount().Currency(),
		tx.Category().Name(),
		string(tx.Category().Type()),
		tx.Description(),
		string(tx.Type()),
		tx.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) FindByID(id string) (transaction.Transaction, error) {
	query := `
		SELECT id, amount_cents, currency, category_name, category_type, description, type, created_at
		FROM transactions
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	return r.scanTransaction(row)
}

func (r *TransactionRepository) FindAll() ([]transaction.Transaction, error) {
	query := `
		SELECT id, amount_cents, currency, category_name, category_type, description, type, created_at
		FROM transactions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) FindByDateRange(start, end time.Time) ([]transaction.Transaction, error) {
	query := `
		SELECT id, amount_cents, currency, category_name, category_type, description, type, created_at
		FROM transactions
		WHERE created_at >= ? AND created_at <= ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) FindByMonth(year int, month time.Month) ([]transaction.Transaction, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	return r.FindByDateRange(start, end)
}

func (r *TransactionRepository) Delete(id string) error {
	query := `DELETE FROM transactions WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return shared.ErrNotFound
	}

	return nil
}

func (r *TransactionRepository) scanTransaction(row *sql.Row) (transaction.Transaction, error) {
	var (
		id           string
		amountCents  int64
		currency     string
		categoryName string
		categoryType string
		description  string
		txType       string
		createdAt    time.Time
	)

	err := row.Scan(
		&id,
		&amountCents,
		&currency,
		&categoryName,
		&categoryType,
		&description,
		&txType,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return transaction.Transaction{}, shared.ErrNotFound
	}

	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("failed to scan transaction: %w", err)
	}

	money := shared.UnsafeFromCents(amountCents)
	category, _ := shared.NewCategory(categoryName, shared.CategoryType(categoryType))

	return transaction.NewTransaction(
		id,
		money,
		category,
		description,
		transaction.TransactionType(txType),
		createdAt,
	), nil
}

func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction

	for rows.Next() {
		var (
			id           string
			amountCents  int64
			currency     string
			categoryName string
			categoryType string
			description  string
			txType       string
			createdAt    time.Time
		)

		err := rows.Scan(
			&id,
			&amountCents,
			&currency,
			&categoryName,
			&categoryType,
			&description,
			&txType,
			&createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		money := shared.UnsafeFromCents(amountCents)
		category, _ := shared.NewCategory(categoryName, shared.CategoryType(categoryType))

		tx := transaction.NewTransaction(
			id,
			money,
			category,
			description,
			transaction.TransactionType(txType),
			createdAt,
		)

		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}
