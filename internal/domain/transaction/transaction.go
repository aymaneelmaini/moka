package transaction

import (
	"moka/internal/shared"
	"time"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	id          string
	amount      shared.Money
	category    shared.Category
	description string
	typ         TransactionType
	createdAt   time.Time
}

func NewTransaction(
	id string,
	amount shared.Money,
	category shared.Category,
	description string,
	typ TransactionType,
	createdAt time.Time,
) Transaction {
	return Transaction{
		id:          id,
		amount:      amount,
		category:    category,
		description: description,
		typ:         typ,
		createdAt:   createdAt,
	}
}

func (t Transaction) ID() string                     { return t.id }
func (t Transaction) Amount() shared.Money           { return t.amount }
func (t Transaction) Category() shared.Category      { return t.category }
func (t Transaction) Description() string            { return t.description }
func (t Transaction) Type() TransactionType          { return t.typ }
func (t Transaction) CreatedAt() time.Time           { return t.createdAt }

func (t Transaction) IsIncome() bool {
	return t.typ == TransactionTypeIncome
}

func (t Transaction) IsExpense() bool {
	return t.typ == TransactionTypeExpense
}
