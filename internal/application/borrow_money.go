package application

import (
	"fmt"
	"moka/internal/domain/loan"
	"moka/internal/domain/transaction"
	"moka/internal/shared"
	"time"

	"github.com/google/uuid"
)

type BorrowMoneyUseCase struct {
	loanRepo        loan.Repository
	transactionRepo transaction.Repository
}

func NewBorrowMoneyUseCase(
	loanRepo loan.Repository,
	transactionRepo transaction.Repository,
) *BorrowMoneyUseCase {
	return &BorrowMoneyUseCase{
		loanRepo:        loanRepo,
		transactionRepo: transactionRepo,
	}
}

type BorrowMoneyInput struct {
	LenderName  string
	Amount      float64
	Description string
	Date        time.Time
}

type BorrowMoneyOutput struct {
	Loan        loan.Loan
	Transaction transaction.Transaction
}

func (uc *BorrowMoneyUseCase) Execute(input BorrowMoneyInput) (*BorrowMoneyOutput, error) {
	// Validate input
	if input.LenderName == "" {
		return nil, fmt.Errorf("lender name cannot be empty: %w", shared.ErrInvalidInput)
	}
	if input.Description == "" {
		return nil, fmt.Errorf("description cannot be empty: %w", shared.ErrInvalidInput)
	}

	money, err := shared.NewMoney(input.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid loan amount: %w", err)
	}

	loanObj := loan.NewLoan(
		uuid.New().String(),
		input.LenderName,
		money,
		input.Date,
		input.Description,
	)

	if err := uc.loanRepo.Save(loanObj); err != nil {
		return nil, fmt.Errorf("failed to save loan: %w", err)
	}

	tx := transaction.NewTransaction(
		uuid.New().String(),
		money,
		shared.CategoryBorrowed,
		fmt.Sprintf("Borrowed from %s: %s", input.LenderName, input.Description),
		transaction.TransactionTypeIncome,
		input.Date,
	)

	if err := uc.transactionRepo.Save(tx); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	return &BorrowMoneyOutput{
		Loan:        loanObj,
		Transaction: tx,
	}, nil
}
