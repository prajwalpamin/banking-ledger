package service

import (
	"context"
	"errors"
	"time"

	"github.com/prajwalpamin/banking-ledger/internal/domain"

	"github.com/google/uuid"
)

type TransactionService struct {
	accountRepo     domain.AccountRepository
	transactionRepo domain.TransactionRepository
	messageProducer domain.MessageProducer
}

func NewTransactionService(
	accountRepo domain.AccountRepository,
	transactionRepo domain.TransactionRepository,
	messageProducer domain.MessageProducer,
) *TransactionService {
	return &TransactionService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		messageProducer: messageProducer,
	}
}

// ProcessMessage implements domain.MessageProcessor
func (s *TransactionService) ProcessMessage(ctx context.Context, msg *domain.Message) (*domain.Transaction, error) {
	return s.ProcessTransaction(msg.AccountID, msg.Type, msg.Amount)
}

func (s *TransactionService) CreateTransaction(ctx context.Context, accountID uuid.UUID, txType domain.TransactionType, amount float64) error {
	// Create message
	msg := &domain.Message{
		AccountID: accountID,
		Type:      txType,
		Amount:    amount,
	}

	// Publish to Kafka
	return s.messageProducer.PublishMessage(ctx, msg)
}

func (s *TransactionService) GetTransactionByID(transactionID uuid.UUID) (*domain.Transaction, error) {
	return &domain.Transaction{}, nil
}

func (s *TransactionService) ProcessTransaction(
	accountID uuid.UUID,
	transactionType domain.TransactionType,
	amount float64,
) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}

	balanceBefore := account.Balance
	var balanceAfter float64

	switch transactionType {
	case domain.TransactionTypeDeposit:
		balanceAfter = balanceBefore + amount
	case domain.TransactionTypeWithdrawal:
		if balanceBefore < amount {
			return nil, errors.New("insufficient funds")
		}
		balanceAfter = balanceBefore - amount
	default:
		return nil, errors.New("invalid transaction type")
	}

	transaction := &domain.Transaction{
		ID:            uuid.New(),
		AccountID:     accountID,
		Type:          transactionType,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        domain.TransactionStatusPending,
		CreatedAt:     time.Now(),
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	account.Balance = balanceAfter
	account.UpdatedAt = time.Now()

	if err := s.accountRepo.Update(account); err != nil {
		transaction.Status = domain.TransactionStatusFailed
		_ = s.transactionRepo.Update(transaction)
		return nil, err
	}

	transaction.Status = domain.TransactionStatusCompleted
	if err := s.transactionRepo.Update(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactionHistory(accountID uuid.UUID) ([]domain.Transaction, error) {
	return s.transactionRepo.GetByAccountID(accountID)
}
