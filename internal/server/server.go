package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/prajwalpamin/banking-ledger/config"
	"github.com/prajwalpamin/banking-ledger/internal/service"
	"github.com/prajwalpamin/banking-ledger/pkg/queue"
)

type Server struct {
	cfg            *config.Config
	router         *gin.Engine
	httpServer     *http.Server
	accountSvc     *service.AccountService
	transactionSvc *service.TransactionService
}

func NewServer(
	cfg *config.Config,
	db *gorm.DB,
	mongoClient *mongo.Client,
	kafkaClient *queue.KafkaClient,
) *Server {
	router := gin.Default()

	// Initialize services
	accountSvc := service.NewAccountService(db)
	transactionSvc := service.NewTransactionService(db, mongoClient, kafkaClient)

	server := &Server{
		cfg:            cfg,
		router:         router,
		accountSvc:     accountSvc,
		transactionSvc: transactionSvc,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		// Account routes
		api.POST("/accounts", s.CreateAccount)
		api.GET("/accounts/:id", s.GetAccount)
		api.GET("/accounts/:id/balance", s.GetBalance)

		// Transaction routes
		api.POST("/transactions", s.CreateTransaction)
		api.GET("/accounts/:id/transactions", s.GetTransactions)
	}
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.ServerPort),
		Handler: s.router,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// ProcessTransaction is used by the Kafka consumer
func (s *Server) ProcessTransaction(msg *queue.TransactionMessage) error {
	return s.transactionSvc.ProcessTransactionMessage(context.Background(), msg)
}
