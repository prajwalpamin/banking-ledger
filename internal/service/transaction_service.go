package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"

	"github.com/prajwalpamin/banking-ledger/internal/models"
	"github.com/prajwalpamin/banking-ledger/pkg/queue"
	"github.com/prajwalpamin/banking-ledger/pkg/utils"
)

type TransactionService struct {
	db          *gorm.DB
	mongoClient *mongo.Client
	kafkaClient *queue.KafkaClient
}

func NewTransactionService(
	db *gorm.DB,
	mongoClient *mongo.Client,
	kafkaClient *queue.KafkaClient,
) *TransactionService {
	return &TransactionService{
		db:          db,
		mongoClient: mongoClient,
		kafkaClient: kafkaClient,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, accountID string, amount float64, txType string) error {
	// Validate transaction type
	if txType != "deposit" && txType != "withdrawal" {
		return errors.New("invalid transaction type")
	}

	// Check if account exists
	var account models.Account
	if err := s.db.WithContext(ctx).First(&account, "id = ?", accountID).Error; err != nil {
		return errors.New("account not found")
	}

	// For withdrawals, check if sufficient balance
	if txType == "withdrawal" {
		if account.Balance < amount {
			return errors.New("insufficient funds")
		}
	}
	generator := utils.GetIDGenerator()
	transaction_id, err := generator.GenerateTransactionID()
	if err != nil {
		return errors.New("Error in transaction id generation")
	}
	// Create transaction message
	msg := &queue.TransactionMessage{
		TransactionID: transaction_id,
		AccountID:     accountID,
		Amount:        amount,
		Type:          txType,
		Timestamp:     time.Now(),
	}

	// Publish to Kafka
	return s.kafkaClient.PublishTransaction(ctx, msg)
}

func (s *TransactionService) ProcessTransactionMessage(ctx context.Context, msg *queue.TransactionMessage) error {
	// Start a database transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate amount based on transaction type
	amount := msg.Amount
	if msg.Type == "withdrawal" {
		amount = -amount
	}

	// Update account balance
	if err := tx.Model(&models.Account{}).
		Where("id = ?", msg.AccountID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Store transaction in MongoDB
	transaction := models.Transaction{
		ID:          msg.TransactionID,
		AccountID:   msg.AccountID,
		Amount:      msg.Amount,
		Type:        msg.Type,
		Status:      "completed",
		CreatedAt:   msg.Timestamp,
		ProcessedAt: time.Now(),
	}

	collection := s.mongoClient.Database("banking").Collection("transactions")
	_, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (s *TransactionService) GetTransactions(ctx context.Context, accountID string, limit int64) ([]models.Transaction, error) {
	collection := s.mongoClient.Database("banking").Collection("transactions")

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{"account_id": accountID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
