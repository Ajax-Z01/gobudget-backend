package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitDatabase initializes the database connection and performs migrations
func InitDatabase() {
	// Construct the Data Source Name (DSN) from environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),     // Database host
		os.Getenv("DB_USER"),     // Database user
		os.Getenv("DB_PASSWORD"), // Database password
		os.Getenv("DB_NAME"),     // Database name
		os.Getenv("DB_PORT"),     // Database port
	)

	var err error
	// Open the database connection using GORM with PostgreSQL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Prevent table names from being pluralized
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err) // Log and exit if connection fails
	}

	// Drop existing tables (for development/testing purposes)
	DB.Migrator().DropTable(&User{}, &Category{}, &Transaction{}, &Budget{})

	// AutoMigrate creates or updates tables based on the struct definitions
	DB.AutoMigrate(&User{}, &Category{}, &Transaction{}, &Budget{})

	fmt.Println("Database connected & migrated successfully!") // Print success message
}
