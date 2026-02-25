package fixed_charge

import "moka/internal/shared"

// CalculateTotalCharges calculates total of all active fixed charges (pure function)
func CalculateTotalCharges(charges []FixedCharge) shared.Money {
	total := shared.Zero()

	for _, charge := range charges {
		if charge.IsActive() {
			total = total.Add(charge.Amount())
		}
	}

	return total
}

// FilterActive returns only active fixed charges
func FilterActive(charges []FixedCharge) []FixedCharge {
	var active []FixedCharge

	for _, charge := range charges {
		if charge.IsActive() {
			active = append(active, charge)
		}
	}

	return active
}
