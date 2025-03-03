package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDatabase()

	r := gin.Default()

	r.GET("/transactions", GetTransactions)
	r.POST("/transactions", CreateTransaction)
	r.GET("/transactions/:id", GetTransactionByID)
	r.PUT("/transactions/:id", UpdateTransaction)
	r.PUT("/transactions/delete/:id", SoftDeleteTransaction)
	r.PUT("/transactions/restore/:id", RestoreTransaction)
	r.GET("/summary", GetSummary)

	r.Run(":8080")
}
