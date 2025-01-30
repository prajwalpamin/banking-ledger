package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountRepository interface {
	Create(account *Account) error
	GetByID(id uuid.UUID) (*Account, error)
	Update(account *Account) error
	Lock(id uuid.UUID) error
	Unlock(id uuid.UUID) error
}

// type Transaction struct {
// 	ID          string    `json:"id"`
// 	AccountID   string    `json:"account_id"`
// 	Type        string    `json:"type"` // deposit or withdrawal
// 	Amount      float64   `json:"amount"`
// 	Status      string    `json:"status"` // pending, completed, failed
// 	CreatedAt   time.Time `json:"created_at"`
// 	ProcessedAt time.Time `json:"processed_at"`
// }
