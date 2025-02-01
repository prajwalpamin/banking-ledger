package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/prajwalpamin/banking-ledger/internal/models"
)

// Request/Response structures
type CreateAccountRequest struct {
	InitialBalance float64 `json:"initial_balance" binding:"required,gte=0"`
}

type CreateTransactionRequest struct {
	AccountID string  `json:"account_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Type      string  `json:"type" binding:"required,oneof=deposit withdrawal"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// CreateAccount handles account creation requests
func (s *Server) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	account, err := s.accountSvc.CreateAccount(c.Request.Context(), req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create account: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount retrieves account details
func (s *Server) GetAccount(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Account ID is required"})
		return
	}

	account, err := s.accountSvc.GetAccount(c.Request.Context(), accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve account: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetBalance retrieves the current balance for an account
func (s *Server) GetBalance(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Account ID is required"})
		return
	}

	balance, err := s.accountSvc.GetBalance(c.Request.Context(), accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve balance: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"balance":    balance,
	})
}

// CreateTransaction initiates a new transaction
func (s *Server) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	err := s.transactionSvc.CreateTransaction(
		c.Request.Context(),
		req.AccountID,
		req.Amount,
		req.Type,
	)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Account not found"})
		case errors.Is(err, models.ErrInsufficientFunds):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Insufficient funds"})
		case errors.Is(err, models.ErrInvalidTransactionType):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid transaction type"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create transaction: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusAccepted, MessageResponse{Message: "Transaction processing"})
}

// GetTransactions retrieves transaction history for an account
func (s *Server) GetTransactions(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Account ID is required"})
		return
	}

	// Parse query parameters
	limit := 50 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid limit parameter"})
			return
		}
		if parsedLimit > 0 {
			limit = int(parsedLimit)
		}
	}

	transactions, err := s.transactionSvc.GetTransactions(c.Request.Context(), accountID, int64(limit))
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve transactions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":   accountID,
		"transactions": transactions,
	})
}
