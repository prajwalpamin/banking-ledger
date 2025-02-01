package models

import (
	"context"
	"time"
)

// internal/models/transaction.go
type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Type        string    `json:"type"` // "deposit" or "withdrawal"
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt time.Time `json:"processed_at"`
}

// internal/service/transaction_service.go
type TransactionService interface {
	ProcessTransaction(ctx context.Context, accountID string, amount float64, txType string) error
	GetTransactionHistory(ctx context.Context, accountID string) ([]Transaction, error)
}
