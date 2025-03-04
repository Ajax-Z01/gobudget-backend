package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDatabase()
	SeedDatabase()

	r := gin.Default()

	r.GET("/transactions", GetTransactions)
	r.POST("/transactions", CreateTransaction)
	r.GET("/transactions/:id", GetTransactionByID)
	r.PUT("/transactions/:id", UpdateTransaction)
	r.PUT("/transactions/delete/:id", SoftDeleteTransaction)
	r.PUT("/transactions/restore/:id", RestoreTransaction)
	r.POST("/categories", CreateCategory)
	r.GET("/categories", GetCategories)
	r.GET("/categories/:id/transactions", GetTransactionsByCategory)
	r.PUT("/transactions/:id/category", UpdateTransactionCategory)
	r.GET("/summary", GetSummary)

	r.Run(":8080")
}
