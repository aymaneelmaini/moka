package loan

// Repository defines the interface for loan persistence (port)
type Repository interface {
	Save(l Loan) error
	FindByID(id string) (Loan, error)
	FindAll() ([]Loan, error)
	FindActive() ([]Loan, error)
	FindByStatus(status LoanStatus) ([]Loan, error)
	Update(l Loan) error
	Delete(id string) error
}
