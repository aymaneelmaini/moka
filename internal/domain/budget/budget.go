package budget

import (
	"moka/internal/shared"
	"time"
)

type Budget struct {
	id           string
	category     shared.Category
	limit        shared.Money
	month        time.Month
	year         int
}

func NewBudget(
	id string,
	category shared.Category,
	limit shared.Money,
	month time.Month,
	year int,
) Budget {
	return Budget{
		id:       id,
		category: category,
		limit:    limit,
		month:    month,
		year:     year,
	}
}

func (b Budget) ID() string                { return b.id }
func (b Budget) Category() shared.Category { return b.category }
func (b Budget) Limit() shared.Money       { return b.limit }
func (b Budget) Month() time.Month         { return b.month }
func (b Budget) Year() int                 { return b.year }

func (b Budget) IsExceeded(spent shared.Money) bool {
	return spent.GreaterThan(b.limit)
}

func (b Budget) RemainingAmount(spent shared.Money) shared.Money {
	return b.limit.Subtract(spent)
}

func (b Budget) PercentageUsed(spent shared.Money) float64 {
	if b.limit.IsZero() {
		return 0
	}

	return (spent.Amount() / b.limit.Amount()) * 100
}
