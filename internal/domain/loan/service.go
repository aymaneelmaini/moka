package loan

import "github.com/aymaneelmaini/moka/internal/shared"

// CalculateTotalOwed calculates total amount still owed across all active loans (pure function)
func CalculateTotalOwed(loans []Loan) shared.Money {
	total := shared.Zero()

	for _, loan := range loans {
		if loan.IsActive() {
			remaining := loan.RemainingAmount()
			total = total.Add(remaining)
		}
	}

	return total
}

// FilterActive returns only active loans (pure function)
func FilterActive(loans []Loan) []Loan {
	var active []Loan

	for _, loan := range loans {
		if loan.IsActive() {
			active = append(active, loan)
		}
	}

	return active
}

// FilterByLender returns loans from a specific lender (pure function)
func FilterByLender(loans []Loan, lenderName string) []Loan {
	var filtered []Loan

	for _, loan := range loans {
		if loan.LenderName() == lenderName {
			filtered = append(filtered, loan)
		}
	}

	return filtered
}
