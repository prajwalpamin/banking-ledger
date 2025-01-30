package service

import (
	"context"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/google/uuid"
	"github.com/prajwalpamin/banking-ledger/internal/domain"
	"github.com/prajwalpamin/banking-ledger/internal/queue/kafka"
)
type AccountService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	kafkaProducer   kafka.Producer
	logger          *logger.Logger
}

func (s *AccountService) CreateAccount(ctx context.Context, initialBalance float64) (*domain.Account, error) {
	account := &domain.Account{
		ID:        uuid.New().String(),
		Balance:   initialBalance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.accountRepo.Create(ctx, account)
}

func (s *AccountService) ProcessTransaction(ctx context.Context, tx *domain.Transaction) error {
	// Produce transaction to Kafka
	return s.kafkaProducer.ProduceMessage(ctx, "transactions", tx)
}
