package loan

import (
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"
)

type LoanStatus string

const (
	LoanStatusActive   LoanStatus = "active"
	LoanStatusPaidBack LoanStatus = "paid_back"
)

type Loan struct {
	id          string
	lenderName  string          // marouane, younes, hamza, soufiane and and and
	amount      shared.Money
	amountPaid  shared.Money
	borrowedAt  time.Time
	paidBackAt  *time.Time
	status      LoanStatus
	description string
}

func NewLoan(
	id string,
	lenderName string,
	amount shared.Money,
	borrowedAt time.Time,
	description string,
) Loan {
	return Loan{
		id:          id,
		lenderName:  lenderName,
		amount:      amount,
		amountPaid:  shared.Zero(),
		borrowedAt:  borrowedAt,
		paidBackAt:  nil,
		status:      LoanStatusActive,
		description: description,
	}
}

func (l Loan) ID() string            { return l.id }
func (l Loan) LenderName() string    { return l.lenderName }
func (l Loan) Amount() shared.Money  { return l.amount }
func (l Loan) AmountPaid() shared.Money { return l.amountPaid }
func (l Loan) BorrowedAt() time.Time { return l.borrowedAt }
func (l Loan) PaidBackAt() *time.Time { return l.paidBackAt }
func (l Loan) Status() LoanStatus    { return l.status }
func (l Loan) Description() string   { return l.description }

func (l Loan) RemainingAmount() shared.Money {
	return l.amount.Subtract(l.amountPaid)
}

func (l Loan) IsFullyPaid() bool {
	return l.status == LoanStatusPaidBack
}

func (l Loan) IsActive() bool {
	return l.status == LoanStatusActive
}

func (l Loan) RecordPayment(payment shared.Money, paidAt time.Time) Loan {
	newAmountPaid := l.amountPaid.Add(payment)
	remaining := l.amount.Subtract(newAmountPaid)

	newStatus := l.status
	var newPaidBackAt *time.Time

	if remaining.IsZero() || remaining.IsNegative() {
		newStatus = LoanStatusPaidBack
		newPaidBackAt = &paidAt
	}

	return Loan{
		id:          l.id,
		lenderName:  l.lenderName,
		amount:      l.amount,
		amountPaid:  newAmountPaid,
		borrowedAt:  l.borrowedAt,
		paidBackAt:  newPaidBackAt,
		status:      newStatus,
		description: l.description,
	}
}
