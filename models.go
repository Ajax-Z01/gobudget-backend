package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Global variable to hold the database connection
var DB *gorm.DB

// User model representing a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `json:"-"` // Exclude password from JSON responses
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Soft delete support
}

// Category model representing a transaction category
type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"` // Ensure unique category names
}

// Budget model representing budget allocations per category
type Budget struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	CategoryID uint           `gorm:"not null" json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Amount     float64        `gorm:"not null" json:"amount"`
	Spent      float64        `gorm:"-" json:"spent"`
	Month      string         `json:"month"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Transaction model representing income and expenses
type Transaction struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Type       string         `gorm:"not null" json:"type"` // "Income" or "Expense"
	Amount     float64        `gorm:"not null" json:"amount"`
	Note       string         `json:"note"` // Optional transaction note
	CategoryID *uint          `json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"` // Establish foreign key
	UserID     uint           `gorm:"not null" json:"user_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Soft delete support
}

// HashPassword hashes the user's password before storing it in the database
func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the provided password against the hashed password
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
