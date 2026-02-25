package application

import (
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/loan"
	"github.com/aymaneelmaini/moka/internal/domain/transaction"
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"

	"github.com/google/uuid"
)

type PayLoanUseCase struct {
	loanRepo        loan.Repository
	transactionRepo transaction.Repository
}

func NewPayLoanUseCase(
	loanRepo loan.Repository,
	transactionRepo transaction.Repository,
) *PayLoanUseCase {
	return &PayLoanUseCase{
		loanRepo:        loanRepo,
		transactionRepo: transactionRepo,
	}
}

type PayLoanInput struct {
	LoanID  string
	Amount  float64
	Date    time.Time
}

type PayLoanOutput struct {
	UpdatedLoan       loan.Loan
	Transaction       transaction.Transaction
	RemainingAmount   shared.Money
	FullyPaid         bool
}

func (uc *PayLoanUseCase) Execute(input PayLoanInput) (*PayLoanOutput, error) {
	// Validate input
	if input.LoanID == "" {
		return nil, fmt.Errorf("loan ID cannot be empty: %w", shared.ErrInvalidInput)
	}

	payment, err := shared.NewMoney(input.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid payment amount: %w", err)
	}

	loanObj, err := uc.loanRepo.FindByID(input.LoanID)
	if err != nil {
		return nil, fmt.Errorf("failed to find loan: %w", err)
	}

	updatedLoan := loanObj.RecordPayment(payment, input.Date)

	if err := uc.loanRepo.Update(updatedLoan); err != nil {
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	category, _ := shared.NewCategory(
		fmt.Sprintf("Loan Payment - %s", updatedLoan.LenderName()),
		shared.CategoryTypeExpense,
	)

	tx := transaction.NewTransaction(
		uuid.New().String(),
		payment,
		category,
		fmt.Sprintf("Paid %s to %s", payment.String(), updatedLoan.LenderName()),
		transaction.TransactionTypeExpense,
		input.Date,
	)

	if err := uc.transactionRepo.Save(tx); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	return &PayLoanOutput{
		UpdatedLoan:     updatedLoan,
		Transaction:     tx,
		RemainingAmount: updatedLoan.RemainingAmount(),
		FullyPaid:       updatedLoan.IsFullyPaid(),
	}, nil
}
