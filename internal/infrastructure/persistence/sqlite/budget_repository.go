package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/budget"
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"
)

type BudgetRepository struct {
	db *DB
}

func NewBudgetRepository(db *DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Save(b budget.Budget) error {
	query := `
		INSERT INTO budgets (id, category_name, limit_amount, currency, month, year)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		b.ID(),
		b.Category().Name(),
		b.Limit().Amount(),
		b.Limit().Currency(),
		int(b.Month()),
		b.Year(),
	)

	if err != nil {
		return fmt.Errorf("failed to save budget: %w", err)
	}

	return nil
}

func (r *BudgetRepository) FindByID(id string) (budget.Budget, error) {
	query := `
		SELECT id, category_name, limit_amount, currency, month, year
		FROM budgets
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	return r.scanBudget(row)
}

func (r *BudgetRepository) FindByMonthAndYear(month time.Month, year int) ([]budget.Budget, error) {
	query := `
		SELECT id, category_name, limit_amount, currency, month, year
		FROM budgets
		WHERE month = ? AND year = ?
	`

	rows, err := r.db.Query(query, int(month), year)
	if err != nil {
		return nil, fmt.Errorf("failed to query budgets: %w", err)
	}
	defer rows.Close()

	return r.scanBudgets(rows)
}

func (r *BudgetRepository) FindByCategoryAndMonth(categoryName string, month time.Month, year int) (budget.Budget, error) {
	query := `
		SELECT id, category_name, limit_amount, currency, month, year
		FROM budgets
		WHERE category_name = ? AND month = ? AND year = ?
	`

	row := r.db.QueryRow(query, categoryName, int(month), year)
	return r.scanBudget(row)
}

func (r *BudgetRepository) Update(b budget.Budget) error {
	query := `
		UPDATE budgets
		SET category_name = ?, limit_amount = ?, currency = ?, month = ?, year = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		b.Category().Name(),
		b.Limit().Amount(),
		b.Limit().Currency(),
		int(b.Month()),
		b.Year(),
		b.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
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

func (r *BudgetRepository) Delete(id string) error {
	query := `DELETE FROM budgets WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
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

func (r *BudgetRepository) scanBudget(row *sql.Row) (budget.Budget, error) {
	var (
		id           string
		categoryName string
		limitAmount  float64
		currency     string
		month        int
		year         int
	)

	err := row.Scan(&id, &categoryName, &limitAmount, &currency, &month, &year)

	if err == sql.ErrNoRows {
		return budget.Budget{}, shared.ErrNotFound
	}

	if err != nil {
		return budget.Budget{}, fmt.Errorf("failed to scan budget: %w", err)
	}

	money := shared.UnsafeNewMoney(limitAmount)
	category, _ := shared.NewCategory(categoryName, shared.CategoryTypeExpense)

	return budget.NewBudget(id, category, money, time.Month(month), year), nil
}

func (r *BudgetRepository) scanBudgets(rows *sql.Rows) ([]budget.Budget, error) {
	var budgets []budget.Budget

	for rows.Next() {
		var (
			id           string
			categoryName string
			limitAmount  float64
			currency     string
			month        int
			year         int
		)

		err := rows.Scan(&id, &categoryName, &limitAmount, &currency, &month, &year)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}

		money := shared.UnsafeNewMoney(limitAmount)
		category, _ := shared.NewCategory(categoryName, shared.CategoryTypeExpense)

		b := budget.NewBudget(id, category, money, time.Month(month), year)
		budgets = append(budgets, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating budgets: %w", err)
	}

	return budgets, nil
}
