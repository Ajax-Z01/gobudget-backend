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
	query := DB.Preload("Category").Where("user_id = ? AND deleted_at IS NULL", userID)

	// Apply filters if provided
	if startDateStr, endDateStr := c.Query("start_date"), c.Query("end_date"); startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 == nil && err2 == nil {
			query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
		}
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

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transaction Transaction
	if err := DB.Preload("Category").Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
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

	DB.Preload("Category").First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, transaction)
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

// Soft delete transaction
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
	if err := DB.Model(&transaction).Update("deleted_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted", "deleted_at": now})
}

// Restore a soft-deleted transaction
func RestoreTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transaction Transaction
	if err := DB.Unscoped().Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := DB.Unscoped().Model(&transaction).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore transaction"})
		return
	}

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
	DB.Preload("Category").Where("category_id = ?", c.Param("id")).Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

// GetSummary returns income, expenses, and balance summary
func GetSummary(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var totalIncome, totalExpense sql.NullFloat64

	// Hitung Total Income
	err1 := DB.Model(&Transaction{}).
		Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, "Income").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome).Error

	// Hitung Total Expense
	err2 := DB.Model(&Transaction{}).
		Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, "Expense").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense).Error

	// Cek jika ada error saat mengambil data income atau expense
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary data"})
		return
	}

	// Struktur untuk menyimpan tren bulanan
	type MonthlySummary struct {
		Month        string  `json:"month"`
		TotalIncome  float64 `json:"total_income"`
		TotalExpense float64 `json:"total_expense"`
	}

	var trends []MonthlySummary

	// Ambil data tren bulanan
	err3 := DB.Model(&Transaction{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') AS month, "+
			"COALESCE(SUM(CASE WHEN type = 'Income' THEN amount ELSE 0 END), 0) AS total_income, "+
			"COALESCE(SUM(CASE WHEN type = 'Expense' THEN amount ELSE 0 END), 0) AS total_expense").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("month").
		Order("month ASC").
		Scan(&trends).Error

	// Cek jika ada error saat mengambil data tren
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trend data"})
		return
	}

	// Kirim respons dengan tren sebagai array kosong jika tidak ada data
	c.JSON(http.StatusOK, gin.H{
		"total_income":  totalIncome.Float64,
		"total_expense": totalExpense.Float64,
		"balance":       totalIncome.Float64 - totalExpense.Float64,
		"trend":         trends, // Jika kosong, akan dikirim sebagai `[]`
	})
}

// GetBudgets retrieves all budgets for the authenticated user
func GetBudgets(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var budgets []Budget
	if err := DB.Preload("Category").Where("user_id = ?", userID).Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budgets"})
		return
	}

	for i := range budgets {
		var totalSpent float64
		DB.Model(&Transaction{}).
			Where("user_id = ? AND category_id = ? AND type = ?", userID, budgets[i].CategoryID, "Expense").
			Select("COALESCE(SUM(amount), 0)").Scan(&totalSpent)
		budgets[i].Spent = totalSpent
	}

	c.JSON(http.StatusOK, budgets)
}

// GetBudgetByID retrieves a single budget by its ID
func GetBudgetByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var budget Budget
	if err := DB.Preload("Category").Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	var totalSpent float64
	DB.Model(&Transaction{}).
		Where("user_id = ? AND category_id = ? AND type = ?", userID, budget.CategoryID, "Expense").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalSpent)

	budget.Spent = totalSpent

	c.JSON(http.StatusOK, budget)
}

// CreateBudget adds a new budget
func CreateBudget(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		CategoryID uint    `json:"category_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget := Budget{
		UserID:     userID.(uint),
		CategoryID: input.CategoryID,
		Amount:     input.Amount,
	}

	if err := DB.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
		return
	}

	DB.Preload("Category").First(&budget, budget.ID)

	c.JSON(http.StatusCreated, budget)
}

// UpdateBudget updates an existing budget
func UpdateBudget(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var budget Budget
	if err := DB.Preload("Category").Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	var input struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget.Amount = input.Amount
	DB.Save(&budget)
	c.JSON(http.StatusOK, budget)
}

// SoftDeleteBudget marks a budget as deleted (soft delete)
func SoftDeleteBudget(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var budget Budget
	if err := DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	if err := DB.Model(&budget).Update("deleted_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted (soft deleted)"})
}

// RestoreBudget restores a soft deleted budget
func RestoreBudget(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var budget Budget
	if err := DB.Unscoped().Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	if err := DB.Model(&budget).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget restored"})
}
