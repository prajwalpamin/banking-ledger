package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID            uuid.UUID         `json:"id" bson:"_id"`
	AccountID     uuid.UUID         `json:"account_id" bson:"account_id"`
	Type          TransactionType   `json:"type" bson:"type"`
	Amount        float64           `json:"amount" bson:"amount"`
	BalanceBefore float64           `json:"balance_before" bson:"balance_before"`
	BalanceAfter  float64           `json:"balance_after" bson:"balance_after"`
	Status        TransactionStatus `json:"status" bson:"status"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByAccountID(accountID uuid.UUID) ([]Transaction, error)
	Update(transaction *Transaction) error
}
