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
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Category model representing a transaction category
type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

// Budget model representing budget allocations per category
type Budget struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	CategoryID   uint           `gorm:"not null" json:"category_id"`
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Amount       float64        `gorm:"not null" json:"amount"`
	Currency     string         `gorm:"not null" json:"currency"`
	ExchangeRate float64        `gorm:"not null" json:"exchange_rate"`         // Exchange Rate to IDR
	Spent        float64        `gorm:"-" json:"spent"`                        // Calculated field
	Month        string         `gorm:"type:varchar(7);not null" json:"month"` // Format: "YYYY-MM"
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Transaction model representing income and expenses
type Transaction struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Type         string         `gorm:"not null" json:"type"` // "Income" or "Expense"
	Amount       float64        `gorm:"not null" json:"amount"`
	Currency     string         `gorm:"not null" json:"currency"`
	ExchangeRate float64        `gorm:"not null" json:"exchange_rate"` // Exchange Rate to IDR
	Note         string         `json:"note"`
	CategoryID   *uint          `json:"category_id"`
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
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
