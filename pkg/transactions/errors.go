package transactions

import "errors"

var (
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrInvalidAmount        = errors.New("invalid amount")
)
