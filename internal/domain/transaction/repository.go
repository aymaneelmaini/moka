package transaction

import "time"

// Repository defines the interface for transaction persistence (port)
type Repository interface {
	Save(tx Transaction) error
	FindByID(id string) (Transaction, error)
	FindAll() ([]Transaction, error)
	FindByDateRange(start, end time.Time) ([]Transaction, error)
	FindByMonth(year int, month time.Month) ([]Transaction, error)
	Delete(id string) error
}
