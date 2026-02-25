package budget

import "time"

// Repository defines the interface for budget persistence (port)
type Repository interface {
	Save(b Budget) error
	FindByID(id string) (Budget, error)
	FindByMonthAndYear(month time.Month, year int) ([]Budget, error)
	FindByCategoryAndMonth(categoryName string, month time.Month, year int) (Budget, error)
	Delete(id string) error
	Update(b Budget) error
}
