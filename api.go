package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	var transactions []Transaction
	DB.Where("deleted_at IS NULL").Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

func CreateTransaction(c *gin.Context) {
	var transaction Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&transaction)
	c.JSON(http.StatusCreated, transaction)
}

func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	var transaction Transaction
	if err := DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

func UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	var transaction Transaction
	if err := DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Save(&transaction)
	c.JSON(http.StatusOK, transaction)
}

func SoftDeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	var transaction Transaction
	if err := DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	now := time.Now()
	DB.Model(&transaction).Update("deleted_at", now)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction flagged as deleted", "deleted_at": now})
}

func RestoreTransaction(c *gin.Context) {
	id := c.Param("id")
	var transaction Transaction
	if err := DB.Unscoped().First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	DB.Model(&transaction).Update("deleted_at", nil)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction restored"})
}

func GetSummary(c *gin.Context) {
	var totalIncome sql.NullFloat64
	var totalExpense sql.NullFloat64

	DB.Model(&Transaction{}).Where("type = ? AND deleted_at IS NULL", "Income").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	DB.Model(&Transaction{}).Where("type = ? AND deleted_at IS NULL", "Expense").Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	c.JSON(http.StatusOK, gin.H{
		"total_income":  totalIncome.Float64,
		"total_expense": totalExpense.Float64,
		"balance":       totalIncome.Float64 - totalExpense.Float64,
	})
}
