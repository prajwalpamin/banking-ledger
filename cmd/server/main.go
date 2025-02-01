package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prajwalpamin/banking-ledger/config"
	"github.com/prajwalpamin/banking-ledger/internal/server"
	"github.com/prajwalpamin/banking-ledger/pkg/database"
	"github.com/prajwalpamin/banking-ledger/pkg/queue"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connections
	db, err := database.NewPostgresConnection(ctx, cfg.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	mongoClient, err := database.NewMongoConnection(ctx, cfg.MongoConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// Initialize Kafka client
	kafkaClient, err := queue.NewKafkaClient(cfg.KafkaConfig.Brokers)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer kafkaClient.Close()

	// Initialize server
	srv := server.NewServer(cfg, db, mongoClient, kafkaClient)

	// Start server
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start Kafka consumer
	go func() {
		if err := kafkaClient.StartConsumer(ctx, srv.ProcessTransaction); err != nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown timeout context
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
