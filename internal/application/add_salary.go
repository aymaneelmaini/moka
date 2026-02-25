package application

import (
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/fixed_charge"
	"github.com/aymaneelmaini/moka/internal/domain/transaction"
	"github.com/aymaneelmaini/moka/internal/shared"
	"time"

	"github.com/google/uuid"
)

type AddSalaryUseCase struct {
	transactionRepo transaction.Repository  
	fixedChargeRepo fixed_charge.Repository
}

func NewAddSalaryUseCase(
	transactionRepo transaction.Repository,
	fixedChargeRepo fixed_charge.Repository,
) *AddSalaryUseCase {
	return &AddSalaryUseCase{
		transactionRepo: transactionRepo,
		fixedChargeRepo: fixedChargeRepo,
	}
}

type AddSalaryInput struct {
	Amount      float64
	Description string
	Date        time.Time
}

type AddSalaryOutput struct {
	SalaryTransaction    transaction.Transaction
	FixedCharges         []fixed_charge.FixedCharge
	FixedChargesTotal    shared.Money
	NetAmount            shared.Money
	ChargeTransactions   []transaction.Transaction
}

func (uc *AddSalaryUseCase) Execute(input AddSalaryInput) (*AddSalaryOutput, error) {
	// Validate input
	if input.Description == "" {
		return nil, shared.ErrInvalidInput
	}

	money, err := shared.NewMoney(input.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid salary amount: %w", err)
	}

	salaryTx := transaction.NewTransaction(
		uuid.New().String(),
		money,
		shared.CategorySalary,
		input.Description,
		transaction.TransactionTypeIncome,
		input.Date,
	)

	if err := uc.transactionRepo.Save(salaryTx); err != nil {
		return nil, fmt.Errorf("failed to save salary transaction: %w", err)
	}

	activeCharges, err := uc.fixedChargeRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed charges: %w", err)
	}

	totalCharges := fixed_charge.CalculateTotalCharges(activeCharges)

	var chargeTransactions []transaction.Transaction
	for _, charge := range activeCharges {
		category, _ := shared.NewCategory(charge.Name(), shared.CategoryTypeExpense)

		chargeTx := transaction.NewTransaction(
			uuid.New().String(),
			charge.Amount(),
			category,
			fmt.Sprintf("Fixed charge: %s", charge.Description()),
			transaction.TransactionTypeExpense,
			input.Date,
		)

		if err := uc.transactionRepo.Save(chargeTx); err != nil {
			return nil, fmt.Errorf("failed to save fixed charge transaction: %w", err)
		}

		chargeTransactions = append(chargeTransactions, chargeTx)
	}

	netAmount := money.Subtract(totalCharges)

	return &AddSalaryOutput{
		SalaryTransaction:  salaryTx,
		FixedCharges:       activeCharges,
		FixedChargesTotal:  totalCharges,
		NetAmount:          netAmount,
		ChargeTransactions: chargeTransactions,
	}, nil
}
