package main

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", Register)
	r.POST("/login", Login)

	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/transactions", GetTransactions)
		auth.POST("/transactions", CreateTransaction)
		auth.GET("/transactions/:id", GetTransactionByID)
		auth.PUT("/transactions/:id", UpdateTransaction)
		auth.PUT("/transactions/delete/:id", SoftDeleteTransaction)
		auth.PUT("/transactions/restore/:id", RestoreTransaction)
		auth.POST("/categories", CreateCategory)
		auth.GET("/categories", GetCategories)
		auth.GET("/categories/:id/transactions", GetTransactionsByCategory)
		auth.PUT("/transactions/:id/category", UpdateTransactionCategory)
		auth.GET("/summary", GetSummary)
	}

	return r
}
