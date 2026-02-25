package application

import (
	"fmt"
	"moka/internal/domain/budget"
	"moka/internal/domain/loan"
	"moka/internal/domain/transaction"
	"moka/internal/shared"
	"time"
)

type GetMonthlySummaryUseCase struct {
	transactionRepo transaction.Repository
	budgetRepo      budget.Repository
	loanRepo        loan.Repository
}

func NewGetMonthlySummaryUseCase(
	transactionRepo transaction.Repository,
	budgetRepo budget.Repository,
	loanRepo loan.Repository,
) *GetMonthlySummaryUseCase {
	return &GetMonthlySummaryUseCase{
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
		loanRepo:        loanRepo,
	}
}

type GetMonthlySummaryInput struct {
	Year  int
	Month time.Month
}

type CategorySummary struct {
	CategoryName   string
	Spent          shared.Money
	Budget         *shared.Money
	Remaining      *shared.Money
	PercentageUsed *float64
	BudgetExceeded bool
}

type GetMonthlySummaryOutput struct {
	Year              int
	Month             time.Month
	TotalIncome       shared.Money
	TotalExpenses     shared.Money
	NetSavings        shared.Money
	Balance           shared.Money
	CategorySummaries []CategorySummary
	TotalLoansOwed    shared.Money
	ActiveLoans       []loan.Loan
	Transactions      []transaction.Transaction
}

func (uc *GetMonthlySummaryUseCase) Execute(input GetMonthlySummaryInput) (*GetMonthlySummaryOutput, error) {
	transactions, err := uc.transactionRepo.FindByMonth(input.Year, input.Month)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	totalIncome := transaction.CalculateMonthlyTotal(transactions, transaction.TransactionTypeIncome)
	totalExpenses := transaction.CalculateMonthlyTotal(transactions, transaction.TransactionTypeExpense)
	netSavings := totalIncome.Subtract(totalExpenses)
	balance := transaction.CalculateBalance(transactions)

	categoryTotals := transaction.CalculateCategoryTotal(transactions)

	budgets, _ := uc.budgetRepo.FindByMonthAndYear(input.Month, input.Year)
	budgetMap := make(map[string]budget.Budget)
	for _, b := range budgets {
		budgetMap[b.Category().Name()] = b
	}

	var categorySummaries []CategorySummary
	for categoryName, spent := range categoryTotals {
		summary := CategorySummary{
			CategoryName: categoryName,
			Spent:        spent,
		}

		if b, exists := budgetMap[categoryName]; exists {
			limit := b.Limit()
			remaining := b.RemainingAmount(spent)
			percentage := b.PercentageUsed(spent)
			exceeded := b.IsExceeded(spent)

			summary.Budget = &limit
			summary.Remaining = &remaining
			summary.PercentageUsed = &percentage
			summary.BudgetExceeded = exceeded
		}

		categorySummaries = append(categorySummaries, summary)
	}

	activeLoans, _ := uc.loanRepo.FindActive()
	totalLoansOwed := loan.CalculateTotalOwed(activeLoans)

	return &GetMonthlySummaryOutput{
		Year:              input.Year,
		Month:             input.Month,
		TotalIncome:       totalIncome,
		TotalExpenses:     totalExpenses,
		NetSavings:        netSavings,
		Balance:           balance,
		CategorySummaries: categorySummaries,
		TotalLoansOwed:    totalLoansOwed,
		ActiveLoans:       activeLoans,
		Transactions:      transactions,
	}, nil
}
