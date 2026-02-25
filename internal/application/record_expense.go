package application

import (
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/budget"
	"github.com/aymaneelmaini/moka/internal/domain/transaction"
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"

	"github.com/google/uuid"
)

type RecordExpenseUseCase struct {
	transactionRepo transaction.Repository
	budgetRepo      budget.Repository
}

func NewRecordExpenseUseCase(
	transactionRepo transaction.Repository,
	budgetRepo budget.Repository,
) *RecordExpenseUseCase {
	return &RecordExpenseUseCase{
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
	}
}

type RecordExpenseInput struct {
	Amount       float64
	CategoryName string
	Description  string
	Date         time.Time
}

type RecordExpenseOutput struct {
	Transaction     transaction.Transaction
	Budget          *budget.Budget
	Spent           shared.Money
	RemainingBudget shared.Money
	BudgetExceeded  bool
	PercentageUsed  float64
}

func (uc *RecordExpenseUseCase) Execute(input RecordExpenseInput) (*RecordExpenseOutput, error) {
	// Validate input
	if input.CategoryName == "" {
		return nil, fmt.Errorf("category name cannot be empty: %w", shared.ErrInvalidInput)
	}
	if input.Description == "" {
		return nil, fmt.Errorf("description cannot be empty: %w", shared.ErrInvalidInput)
	}

	category, err := shared.NewCategory(input.CategoryName, shared.CategoryTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	money, err := shared.NewMoney(input.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid expense amount: %w", err)
	}

	tx := transaction.NewTransaction(
		uuid.New().String(),
		money,
		category,
		input.Description,
		transaction.TransactionTypeExpense,
		input.Date,
	)

	if err := uc.transactionRepo.Save(tx); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	year, month, _ := input.Date.Date()
	budgetObj, err := uc.budgetRepo.FindByCategoryAndMonth(input.CategoryName, month, year)

	output := &RecordExpenseOutput{
		Transaction: tx,
	}

	if err == nil {
		startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

		monthTransactions, err := uc.transactionRepo.FindByDateRange(startOfMonth, endOfMonth)
		if err == nil {
			var categoryTransactions []transaction.Transaction
			for _, t := range monthTransactions {
				if t.Category().Name() == input.CategoryName && t.IsExpense() {
					categoryTransactions = append(categoryTransactions, t)
				}
			}

			spent := transaction.CalculateMonthlyTotal(categoryTransactions, transaction.TransactionTypeExpense)

			output.Budget = &budgetObj
			output.Spent = spent
			output.RemainingBudget = budgetObj.RemainingAmount(spent)
			output.BudgetExceeded = budgetObj.IsExceeded(spent)
			output.PercentageUsed = budgetObj.PercentageUsed(spent)
		}
	}

	return output, nil
}
