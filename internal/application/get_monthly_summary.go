package application

import (
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/budget"
	"github.com/aymaneelmaini/moka/internal/domain/fixed_charge"
	"github.com/aymaneelmaini/moka/internal/domain/loan"
	"github.com/aymaneelmaini/moka/internal/domain/transaction"
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"
)

type GetMonthlySummaryUseCase struct {
	transactionRepo transaction.Repository
	budgetRepo      budget.Repository
	loanRepo        loan.Repository
	fixedChargeRepo fixed_charge.Repository
}

func NewGetMonthlySummaryUseCase(
	transactionRepo transaction.Repository,
	budgetRepo budget.Repository,
	loanRepo loan.Repository,
	fixedChargeRepo fixed_charge.Repository,
) *GetMonthlySummaryUseCase {
	return &GetMonthlySummaryUseCase{
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
		loanRepo:        loanRepo,
		fixedChargeRepo: fixedChargeRepo,
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
	FixedCharges      []fixed_charge.FixedCharge
	FixedChargesTotal shared.Money
	Transactions      []transaction.Transaction
}

func (uc *GetMonthlySummaryUseCase) Execute(input GetMonthlySummaryInput) (*GetMonthlySummaryOutput, error) {
	transactions, err := uc.transactionRepo.FindByMonth(input.Year, input.Month)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	totalIncome := transaction.CalculateMonthlyTotal(transactions, transaction.TransactionTypeIncome)
	totalExpenses := transaction.CalculateMonthlyTotal(transactions, transaction.TransactionTypeExpense)
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

	fixedCharges, _ := uc.fixedChargeRepo.FindActive()
	fixedChargesTotal := fixed_charge.CalculateTotalCharges(fixedCharges)

	totalExpensesWithFixed := totalExpenses.Add(fixedChargesTotal)
	netSavingsWithFixed := totalIncome.Subtract(totalExpensesWithFixed)
	balanceWithFixed := balance.Subtract(fixedChargesTotal)

	return &GetMonthlySummaryOutput{
		Year:              input.Year,
		Month:             input.Month,
		TotalIncome:       totalIncome,
		TotalExpenses:     totalExpensesWithFixed,
		NetSavings:        netSavingsWithFixed,
		Balance:           balanceWithFixed,
		CategorySummaries: categorySummaries,
		TotalLoansOwed:    totalLoansOwed,
		ActiveLoans:       activeLoans,
		FixedCharges:      fixedCharges,
		FixedChargesTotal: fixedChargesTotal,
		Transactions:      transactions,
	}, nil
}
