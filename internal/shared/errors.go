package shared

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInsufficientFund = errors.New("insufficient funds")
	ErrDuplicateEntry   = errors.New("duplicate entry")
)
