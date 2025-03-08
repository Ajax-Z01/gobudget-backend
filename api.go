package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTransactions retrieves transactions for the authenticated user
func GetTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transactions []Transaction
	query := DB.Preload("Category").Preload("User").Where("user_id = ? AND deleted_at IS NULL", userID)

	// Apply filters if provided
	if startDate, endDate := c.Query("start_date"), c.Query("end_date"); startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if txType := c.Query("type"); txType != "" {
		query = query.Where("type = ?", txType)
	}

	query.Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

// CreateTransaction handles adding a new transaction
func CreateTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Type       string  `json:"type" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		Note       string  `json:"note"`
		CategoryID uint    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := Transaction{
		Type:       input.Type,
		Amount:     input.Amount,
		Note:       input.Note,
		CategoryID: &input.CategoryID,
		UserID:     userID.(uint),
	}

	if err := DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	DB.Preload("Category").Preload("User").First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transaction Transaction
	if err := DB.Preload("Category").Preload("User").Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transaction Transaction
	if err := DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var input struct {
		Type       string  `json:"type"`
		Amount     float64 `json:"amount"`
		Note       string  `json:"note"`
		CategoryID uint    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.Type = input.Type
	transaction.Amount = input.Amount
	transaction.Note = input.Note
	transaction.CategoryID = &input.CategoryID

	if err := DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// SoftDeleteTransaction marks a transaction as deleted
func SoftDeleteTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transaction Transaction
	if err := DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	now := time.Now()
	DB.Model(&transaction).Update("deleted_at", now)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted", "deleted_at": now})
}

// RestoreTransaction restores a soft-deleted transaction
func RestoreTransaction(c *gin.Context) {
	var transaction Transaction
	if err := DB.Unscoped().First(&transaction, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	DB.Model(&transaction).Update("deleted_at", nil)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction restored"})
}

// CreateCategory handles adding a new category
func CreateCategory(c *gin.Context) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&category)
	c.JSON(http.StatusCreated, category)
}

// GetCategories retrieves all categories
func GetCategories(c *gin.Context) {
	var categories []Category
	DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

// GetTransactionsByCategory retrieves transactions filtered by category
func GetTransactionsByCategory(c *gin.Context) {
	var transactions []Transaction
	DB.Preload("Category").Preload("User").Where("category_id = ?", c.Param("id")).Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

// UpdateTransactionCategory updates the category of a transaction
func UpdateTransactionCategory(c *gin.Context) {
	var transaction Transaction
	if err := DB.First(&transaction, c.Param("id")).Error; err != nil {
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

// GetSummary returns income, expenses, and balance summary
func GetSummary(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var totalIncome, totalExpense sql.NullFloat64
	DB.Model(&Transaction{}).Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, "Income").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	DB.Model(&Transaction{}).Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, "Expense").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	c.JSON(http.StatusOK, gin.H{
		"total_income":  totalIncome.Float64,
		"total_expense": totalExpense.Float64,
		"balance":       totalIncome.Float64 - totalExpense.Float64,
	})
}
