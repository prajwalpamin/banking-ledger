package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/prajwalpamin/banking-ledger/internal/models"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(ctx context.Context, initialBalance float64) (*models.Account, error) {
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	account := &models.Account{
		Balance:   initialBalance,
		CreatedAt: time.Now(),
	}

	err := s.db.WithContext(ctx).Create(account).Error
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account
	err := s.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	return &account, nil
}

func (s *AccountService) UpdateBalance(ctx context.Context, id string, amount float64) error {
	result := s.db.WithContext(ctx).Model(&models.Account{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}

func (s *AccountService) GetBalance(ctx context.Context, id string) (float64, error) {
	var account models.Account
	err := s.db.WithContext(ctx).Select("balance").First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("account not found")
		}
		return 0, err
	}

	return account.Balance, nil
}
