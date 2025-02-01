package mongodb

import (
	"github.com/prajwalpamin/banking-ledger/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(client *mongo.Client, dbName string) *TransactionRepository {
	collection := database.GetMongoCollection(client, dbName, "transactions")
	return &TransactionRepository{collection: collection}
}