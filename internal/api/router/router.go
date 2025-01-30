package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prajwalpamin/banking-ledger/internal/api/handler"
	"github.com/prajwalpamin/banking-ledger/internal/api/middleware"
)

func SetupRouter(
	accountHandler *handler.AccountHandler,
	transactionHandler *handler.TransactionHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	v1 := router.Group("/api/v1")
	{
		// Account routes
		v1.POST("/accounts", accountHandler.CreateAccount)
		v1.GET("/accounts/:id", accountHandler.GetAccount)

		// Transaction routes
		v1.POST("/transactions", transactionHandler.CreateTransaction)
		v1.GET("/accounts/:id/transactions", transactionHandler.GetTransactionHistory)
	}

	return router
}
