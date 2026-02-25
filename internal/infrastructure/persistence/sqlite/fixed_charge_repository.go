package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/aymaneelmaini/moka/internal/domain/fixed_charge"
	"github.com/aymaneelmaini/moka/internal/shared"
)

type FixedChargeRepository struct {
	db *DB
}

func NewFixedChargeRepository(db *DB) *FixedChargeRepository {
	return &FixedChargeRepository{db: db}
}

func (r *FixedChargeRepository) Save(fc fixed_charge.FixedCharge) error {
	query := `
		INSERT INTO fixed_charges (id, name, amount, currency, description, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		fc.ID(),
		fc.Name(),
		fc.Amount().Amount(),
		fc.Amount().Currency(),
		fc.Description(),
		fc.IsActive(),
	)

	if err != nil {
		return fmt.Errorf("failed to save fixed charge: %w", err)
	}

	return nil
}

func (r *FixedChargeRepository) FindByID(id string) (fixed_charge.FixedCharge, error) {
	query := `
		SELECT id, name, amount, currency, description, is_active
		FROM fixed_charges
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	return r.scanFixedCharge(row)
}

func (r *FixedChargeRepository) FindAll() ([]fixed_charge.FixedCharge, error) {
	query := `
		SELECT id, name, amount, currency, description, is_active
		FROM fixed_charges
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query fixed charges: %w", err)
	}
	defer rows.Close()

	return r.scanFixedCharges(rows)
}

func (r *FixedChargeRepository) FindActive() ([]fixed_charge.FixedCharge, error) {
	query := `
		SELECT id, name, amount, currency, description, is_active
		FROM fixed_charges
		WHERE is_active = 1
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active fixed charges: %w", err)
	}
	defer rows.Close()

	return r.scanFixedCharges(rows)
}

func (r *FixedChargeRepository) Update(fc fixed_charge.FixedCharge) error {
	query := `
		UPDATE fixed_charges
		SET name = ?, amount = ?, currency = ?, description = ?, is_active = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		fc.Name(),
		fc.Amount().Amount(),
		fc.Amount().Currency(),
		fc.Description(),
		fc.IsActive(),
		fc.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update fixed charge: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return shared.ErrNotFound
	}

	return nil
}

func (r *FixedChargeRepository) Delete(id string) error {
	query := `DELETE FROM fixed_charges WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete fixed charge: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// scanFixedCharge scans a single fixed charge from a row
func (r *FixedChargeRepository) scanFixedCharge(row *sql.Row) (fixed_charge.FixedCharge, error) {
	var (
		id          string
		name        string
		amount      float64
		currency    string
		description string
		isActive    bool
	)

	err := row.Scan(&id, &name, &amount, &currency, &description, &isActive)

	if err == sql.ErrNoRows {
		return fixed_charge.FixedCharge{}, shared.ErrNotFound
	}

	if err != nil {
		return fixed_charge.FixedCharge{}, fmt.Errorf("failed to scan fixed charge: %w", err)
	}

	money := shared.UnsafeNewMoney(amount)

	return fixed_charge.NewFixedCharge(id, name, money, description, isActive), nil
}

func (r *FixedChargeRepository) scanFixedCharges(rows *sql.Rows) ([]fixed_charge.FixedCharge, error) {
	var charges []fixed_charge.FixedCharge

	for rows.Next() {
		var (
			id          string
			name        string
			amount      float64
			currency    string
			description string
			isActive    bool
		)

		err := rows.Scan(&id, &name, &amount, &currency, &description, &isActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fixed charge: %w", err)
		}

		money := shared.UnsafeNewMoney(amount)
		fc := fixed_charge.NewFixedCharge(id, name, money, description, isActive)
		charges = append(charges, fc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fixed charges: %w", err)
	}

	return charges, nil
}
