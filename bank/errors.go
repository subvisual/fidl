package bank

import "errors"

var (
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrOperationNotAllowed = errors.New("operation not allowed")
	ErrNothingToRefund     = errors.New("nothing to refund")
	ErrAuthNotFound        = errors.New("authorization not found")
	ErrAuthLocked          = errors.New("authorization is locked")
)
