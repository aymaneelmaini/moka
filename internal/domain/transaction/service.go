package transaction

import (
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"
)

// CalculateBalance calculates total balance from transactions (pure function)
func CalculateBalance(transactions []Transaction) shared.Money {
	balance := shared.Zero()

	for _, tx := range transactions {
		if tx.IsIncome() {
			balance = balance.Add(tx.Amount())
		} else {
			balance = balance.Subtract(tx.Amount())
		}
	}

	return balance
}

// CalculateMonthlyTotal calculates total for a specific month (pure function)
func CalculateMonthlyTotal(transactions []Transaction, typ TransactionType) shared.Money {
	total := shared.Zero()

	for _, tx := range transactions {
		if tx.Type() == typ {
			total = total.Add(tx.Amount())
		}
	}

	return total
}

// GroupByCategory groups transactions by category (pure function)
func GroupByCategory(transactions []Transaction) map[string][]Transaction {
	grouped := make(map[string][]Transaction)

	for _, tx := range transactions {
		categoryName := tx.Category().Name()
		grouped[categoryName] = append(grouped[categoryName], tx)
	}

	return grouped
}

// CalculateCategoryTotal calculates total spent per category (pure function)
func CalculateCategoryTotal(transactions []Transaction) map[string]shared.Money {
	totals := make(map[string]shared.Money)

	for _, tx := range transactions {
		if tx.IsExpense() {
			categoryName := tx.Category().Name()
			current, exists := totals[categoryName]
			if !exists {
				current = shared.Zero()
			}
			totals[categoryName] = current.Add(tx.Amount())
		}
	}

	return totals
}

// FilterByDateRange filters transactions within date range (pure function)
func FilterByDateRange(transactions []Transaction, start, end time.Time) []Transaction {
	var filtered []Transaction

	for _, tx := range transactions {
		if (tx.CreatedAt().After(start) || tx.CreatedAt().Equal(start)) &&
			(tx.CreatedAt().Before(end) || tx.CreatedAt().Equal(end)) {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}
