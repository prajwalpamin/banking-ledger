package models

import (
	"context"
	"time"
)

// internal/models/account.go
type Account struct {
	ID        string    `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// internal/service/account_service.go
type AccountService interface {
	Create(ctx context.Context, initialBalance float64) (*Account, error)
	GetBalance(ctx context.Context, accountID string) (float64, error)
	UpdateBalance(ctx context.Context, accountID string, amount float64) error
}
