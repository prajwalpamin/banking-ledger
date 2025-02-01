package models

import "errors"

var (
	ErrAccountNotFound        = errors.New("account not found")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)
