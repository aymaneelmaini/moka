package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/loan"
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"
)

type LoanRepository struct {
	db *DB
}

func NewLoanRepository(db *DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Save(l loan.Loan) error {
	query := `
		INSERT INTO loans (id, lender_name, amount, amount_paid, currency, borrowed_at, paid_back_at, status, description)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var paidBackAt *time.Time
	if l.PaidBackAt() != nil {
		paidBackAt = l.PaidBackAt()
	}

	_, err := r.db.Exec(
		query,
		l.ID(),
		l.LenderName(),
		l.Amount().Amount(),
		l.AmountPaid().Amount(),
		l.Amount().Currency(),
		l.BorrowedAt(),
		paidBackAt,
		string(l.Status()),
		l.Description(),
	)

	if err != nil {
		return fmt.Errorf("failed to save loan: %w", err)
	}

	return nil
}

func (r *LoanRepository) FindByID(id string) (loan.Loan, error) {
	query := `
		SELECT id, lender_name, amount, amount_paid, currency, borrowed_at, paid_back_at, status, description
		FROM loans
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	return r.scanLoan(row)
}

func (r *LoanRepository) FindAll() ([]loan.Loan, error) {
	query := `
		SELECT id, lender_name, amount, amount_paid, currency, borrowed_at, paid_back_at, status, description
		FROM loans
		ORDER BY borrowed_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query loans: %w", err)
	}
	defer rows.Close()

	return r.scanLoans(rows)
}

func (r *LoanRepository) FindActive() ([]loan.Loan, error) {
	query := `
		SELECT id, lender_name, amount, amount_paid, currency, borrowed_at, paid_back_at, status, description
		FROM loans
		WHERE status = 'active'
		ORDER BY borrowed_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active loans: %w", err)
	}
	defer rows.Close()

	return r.scanLoans(rows)
}

func (r *LoanRepository) FindByStatus(status loan.LoanStatus) ([]loan.Loan, error) {
	query := `
		SELECT id, lender_name, amount, amount_paid, currency, borrowed_at, paid_back_at, status, description
		FROM loans
		WHERE status = ?
		ORDER BY borrowed_at DESC
	`

	rows, err := r.db.Query(query, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to query loans by status: %w", err)
	}
	defer rows.Close()

	return r.scanLoans(rows)
}

func (r *LoanRepository) Update(l loan.Loan) error {
	query := `
		UPDATE loans
		SET lender_name = ?, amount = ?, amount_paid = ?, currency = ?, borrowed_at = ?, paid_back_at = ?, status = ?, description = ?
		WHERE id = ?
	`

	var paidBackAt *time.Time
	if l.PaidBackAt() != nil {
		paidBackAt = l.PaidBackAt()
	}

	result, err := r.db.Exec(
		query,
		l.LenderName(),
		l.Amount().Amount(),
		l.AmountPaid().Amount(),
		l.Amount().Currency(),
		l.BorrowedAt(),
		paidBackAt,
		string(l.Status()),
		l.Description(),
		l.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
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

func (r *LoanRepository) Delete(id string) error {
	query := `DELETE FROM loans WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete loan: %w", err)
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

func (r *LoanRepository) scanLoan(row *sql.Row) (loan.Loan, error) {
	var (
		id            string
		lenderName    string
		amount        float64
		amountPaid    float64
		currency      string
		borrowedAt    time.Time
		paidBackAt    sql.NullTime
		status        string
		description   string
	)

	err := row.Scan(
		&id,
		&lenderName,
		&amount,
		&amountPaid,
		&currency,
		&borrowedAt,
		&paidBackAt,
		&status,
		&description,
	)

	if err == sql.ErrNoRows {
		return loan.Loan{}, shared.ErrNotFound
	}

	if err != nil {
		return loan.Loan{}, fmt.Errorf("failed to scan loan: %w", err)
	}

	amountMoney := shared.UnsafeNewMoney(amount)
	amountPaidMoney := shared.UnsafeNewMoney(amountPaid)

	l := loan.NewLoan(id, lenderName, amountMoney, borrowedAt, description)

	if amountPaid > 0 {
		paymentAmount := amountPaidMoney
		paymentTime := borrowedAt
		if paidBackAt.Valid {
			paymentTime = paidBackAt.Time
		}
		l = l.RecordPayment(paymentAmount, paymentTime)
	}

	return l, nil
}

func (r *LoanRepository) scanLoans(rows *sql.Rows) ([]loan.Loan, error) {
	var loans []loan.Loan

	for rows.Next() {
		var (
			id            string
			lenderName    string
			amount        float64
			amountPaid    float64
			currency      string
			borrowedAt    time.Time
			paidBackAt    sql.NullTime
			status        string
			description   string
		)

		err := rows.Scan(
			&id,
			&lenderName,
			&amount,
			&amountPaid,
			&currency,
			&borrowedAt,
			&paidBackAt,
			&status,
			&description,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan loan: %w", err)
		}

		amountMoney := shared.UnsafeNewMoney(amount)
		amountPaidMoney := shared.UnsafeNewMoney(amountPaid)

		l := loan.NewLoan(id, lenderName, amountMoney, borrowedAt, description)

		if amountPaid > 0 {
			paymentAmount := amountPaidMoney
			paymentTime := borrowedAt
			if paidBackAt.Valid {
				paymentTime = paidBackAt.Time
			}
			l = l.RecordPayment(paymentAmount, paymentTime)
		}

		loans = append(loans, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating loans: %w", err)
	}

	return loans, nil
}
