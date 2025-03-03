package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type Transaction struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Type      string     `json:"type"`
	Amount    float64    `json:"amount"`
	Note      string     `json:"note"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func InitDatabase() {
	dsn := "host=localhost user=postgres password=24'K>W6tMr5an!? dbname=gobudget port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&Transaction{})
	fmt.Println("Database connected & migrated successfully!")
}
