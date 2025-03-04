package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	var transactions []Transaction
	query := DB.Preload("Category").Where("deleted_at IS NULL").Find(&transactions)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	categoryID := c.Query("category_id")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	txType := c.Query("type")
	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	query.Find(&transactions)
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
	if err := DB.Preload("Category").First(&transaction, id).Error; err != nil {
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

func CreateCategory(c *gin.Context) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&category)
	c.JSON(http.StatusCreated, category)
}

func GetCategories(c *gin.Context) {
	var categories []Category
	DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func GetTransactionsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var transactions []Transaction
	DB.Preload("Category").Where("category_id = ?", categoryID).Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

func UpdateTransactionCategory(c *gin.Context) {
	var transaction Transaction
	id := c.Param("id")

	if err := DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var input struct {
		CategoryID uint `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.CategoryID = &input.CategoryID
	DB.Save(&transaction)
	c.JSON(http.StatusOK, transaction)
}

func GetSummary(c *gin.Context) {
	var totalIncome sql.NullFloat64
	var totalExpense sql.NullFloat64

	DB.Model(&Transaction{}).Where("type = ? AND deleted_at IS NULL", "Income").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)

	DB.Model(&Transaction{}).Where("type = ? AND deleted_at IS NULL", "Expense").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	c.JSON(http.StatusOK, gin.H{
		"total_income":  totalIncome.Float64,
		"total_expense": totalExpense.Float64,
		"balance":       totalIncome.Float64 - totalExpense.Float64,
	})
}
