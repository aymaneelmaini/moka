package fixed_charge

// Repository defines the interface for fixed charge persistence
type Repository interface {
	Save(fc FixedCharge) error
	FindByID(id string) (FixedCharge, error)
	FindAll() ([]FixedCharge, error)
	FindActive() ([]FixedCharge, error)
	Update(fc FixedCharge) error
	Delete(id string) error
}
