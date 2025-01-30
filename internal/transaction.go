package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"

	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID            uuid.UUID         `json:"id"`
	AccountID     uuid.UUID         `json:"account_id"`
	Type          TransactionType   `json:"type"`
	Amount        float64           `json:"amount"`
	BalanceBefore float64           `json:"balance_before"`
	BalanceAfter  float64           `json:"balance_after"`
	Status        TransactionStatus `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByAccountID(accountID uuid.UUID) ([]Transaction, error)
	UpdateStatus(id uuid.UUID, status TransactionStatus) error
}
