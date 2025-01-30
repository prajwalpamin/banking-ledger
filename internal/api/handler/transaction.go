package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prajwalpamin/banking-ledger/internal/domain"
	"github.com/prajwalpamin/banking-ledger/internal/service"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

type CreateTransactionRequest struct {
	AccountID uuid.UUID              `json:"account_id" binding:"required"`
	Type      domain.TransactionType `json:"type" binding:"required"`
	Amount    float64                `json:"amount" binding:"required"`
}

type TransactionResponse struct {
	ID            uuid.UUID                `json:"id"`
	AccountID     uuid.UUID                `json:"account_id"`
	Type          domain.TransactionType   `json:"type"`
	Amount        float64                  `json:"amount"`
	BalanceBefore float64                  `json:"balance_before"`
	BalanceAfter  float64                  `json:"balance_after"`
	Status        domain.TransactionStatus `json:"status"`
	CreatedAt     string                   `json:"created_at"`
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Amount must be greater than zero",
		})
		return
	}

	// Validate transaction type
	if req.Type != domain.TransactionTypeDeposit && req.Type != domain.TransactionTypeWithdrawal {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction type. Must be either DEPOSIT or WITHDRAWAL",
		})
		return
	}

	transaction, err := h.transactionService.ProcessTransaction(req.AccountID, req.Type, req.Amount)
	if err != nil {
		// Handle different types of errors
		switch err.Error() {
		case "account not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient funds":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
		}
		return
	}

	response := TransactionResponse{
		ID:            transaction.ID,
		AccountID:     transaction.AccountID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
		Status:        transaction.Status,
		CreatedAt:     transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID format",
		})
		return
	}

	// Get query parameters for pagination
	limit := 10 // default limit
	offset := 0 // default offset
	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if _, err := fmt.Sscanf(offsetParam, "%d", &offset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
	}

	transactions, err := h.transactionService.GetTransactionHistory(accountID)
	if err != nil {
		switch err.Error() {
		case "account not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve transaction history",
			})
		}
		return
	}

	// Convert transactions to response format
	var response []TransactionResponse
	for _, tx := range transactions {
		response = append(response, TransactionResponse{
			ID:            tx.ID,
			AccountID:     tx.AccountID,
			Type:          tx.Type,
			Amount:        tx.Amount,
			BalanceBefore: tx.BalanceBefore,
			BalanceAfter:  tx.BalanceAfter,
			Status:        tx.Status,
			CreatedAt:     tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(response) {
		start = len(response)
	}
	if end > len(response) {
		end = len(response)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":        len(response),
		"limit":        limit,
		"offset":       offset,
		"transactions": response[start:end],
	})
}

// Additional helper methods for the handler

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID format",
		})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(transactionID)
	if err != nil {
		switch err.Error() {
		case "transaction not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve transaction",
			})
		}
		return
	}

	response := TransactionResponse{
		ID:            transaction.ID,
		AccountID:     transaction.AccountID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
		Status:        transaction.Status,
		CreatedAt:     transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}
