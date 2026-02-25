package fixed_charge

import (
	"github.com/aymaneelmaini/moka/internal/shared"
)

type FixedCharge struct {
	id          string
	name        string
	amount      shared.Money
	description string
	isActive    bool
}

func NewFixedCharge(
	id string,
	name string,
	amount shared.Money,
	description string,
	isActive bool,
) FixedCharge {
	return FixedCharge{
		id:          id,
		name:        name,
		amount:      amount,
		description: description,
		isActive:    isActive,
	}
}

func (f FixedCharge) ID() string           { return f.id }
func (f FixedCharge) Name() string         { return f.name }
func (f FixedCharge) Amount() shared.Money { return f.amount }
func (f FixedCharge) Description() string  { return f.description }
func (f FixedCharge) IsActive() bool       { return f.isActive }

func (f FixedCharge) Deactivate() FixedCharge {
	return FixedCharge{
		id:          f.id,
		name:        f.name,
		amount:      f.amount,
		description: f.description,
		isActive:    false,
	}
}

func (f FixedCharge) Activate() FixedCharge {
	return FixedCharge{
		id:          f.id,
		name:        f.name,
		amount:      f.amount,
		description: f.description,
		isActive:    true,
	}
}
