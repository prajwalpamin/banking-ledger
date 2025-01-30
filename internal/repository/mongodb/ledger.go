package mongodb

import (
	"context"
	"time"

	"github.com/prajwalpamin/banking-ledger/internal/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(collection *mongo.Collection) domain.TransactionRepository {
	return &transactionRepository{collection: collection}
}

func (r *transactionRepository) Create(transaction *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, transaction)
	return err
}

func (r *transactionRepository) GetByAccountID(accountID uuid.UUID) ([]domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"account_id": accountID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []domain.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) Update(transaction *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": transaction.ID}
	update := bson.M{"$set": transaction}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
