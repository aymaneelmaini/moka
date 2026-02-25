package shared

import (
	"errors"
	"fmt"
)

const (
	CurrencyMAD = "MAD"
)

var (
	ErrNegativeAmount = errors.New("amount cannot be negative")
	ErrZeroAmount     = errors.New("amount cannot be zero")
	ErrInvalidAmount  = errors.New("amount is invalid")
)

type Money struct {
	amount int64
}

func NewMoney(amount float64) (Money, error) {
	if amount <= 0 {
		if amount == 0 {
			return Money{}, ErrZeroAmount
		}
		return Money{}, ErrNegativeAmount
	}

	cents := int64(amount * 100)
	return Money{amount: cents}, nil
}

func FromCents(cents int64) (Money, error) {
	if cents <= 0 {
		if cents == 0 {
			return Money{}, ErrZeroAmount
		}
		return Money{}, ErrNegativeAmount
	}
	return Money{amount: cents}, nil
}

// UnsafeFromCents creates Money without validation (for internal use like repositories)
func UnsafeFromCents(cents int64) Money {
	return Money{amount: cents}
}

func (m Money) Amount() float64 {
	return float64(m.amount) / 100.0
}

func (m Money) Currency() string {
	return CurrencyMAD
}

func (m Money) AmountInCents() int64 {
	return m.amount
}

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) Subtract(other Money) Money {
	return Money{amount: m.amount - other.amount}
}

func (m Money) IsPositive() bool {
	return m.amount > 0
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) GreaterThan(other Money) bool {
	return m.amount > other.amount
}

func (m Money) LessThan(other Money) bool {
	return m.amount < other.amount
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount(), CurrencyMAD)
}

func Zero() Money {
	return Money{amount: 0}
}
