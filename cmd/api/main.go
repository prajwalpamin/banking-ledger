package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/prajwalpamin/banking-ledger/internal/api/handler"
	"github.com/prajwalpamin/banking-ledger/internal/api/router"
	"github.com/prajwalpamin/banking-ledger/internal/queue/kafka"
	"github.com/prajwalpamin/banking-ledger/internal/repository/mongodb"
	"github.com/prajwalpamin/banking-ledger/internal/repository/postgres"
	"github.com/prajwalpamin/banking-ledger/internal/service"
	"github.com/prajwalpamin/banking-ledger/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initializeDatabases(cfg *config.Config) (*sql.DB, *mongo.Client, error) {
	// Connect to PostgreSQL
	postgresDB, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		return nil, nil, err
	}

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		postgresDB.Close()
		return nil, nil, err
	}

	return postgresDB, mongoClient, nil
}

func initializeKafka(cfg *config.Config, transactionService *service.TransactionService) (*kafka.Producer, *kafka.Consumer, error) {
	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		return nil, nil, err
	}

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		transactionService,
	)
	if err != nil {
		producer.Close()
		return nil, nil, err
	}

	return producer, consumer, nil
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize databases
	postgresDB, mongoClient, err := initializeDatabases(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer postgresDB.Close()
	defer mongoClient.Disconnect(context.Background())

	// Initialize repositories
	accountRepo := postgres.NewAccountRepository(postgresDB)
	transactionRepo := mongodb.NewTransactionRepository(
		mongoClient.Database("banking").Collection("transactions"),
	)

	// Initialize services
	accountService := service.NewAccountService(accountRepo)

	// Create a new transaction service without Kafka producer first
	transactionService := service.NewTransactionService(accountRepo, transactionRepo, nil)

	// Initialize Kafka after service creation
	producer, consumer, err := initializeKafka(cfg, transactionService)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	defer producer.Close()
	defer consumer.Close()

	// Update transaction service with the producer
	transactionService.SetProducer(producer)

	// Start consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	// Initialize handlers
	accountHandler := handler.NewAccountHandler(accountService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Setup router
	r := router.SetupRouter(accountHandler, transactionHandler)

	// Start server
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
