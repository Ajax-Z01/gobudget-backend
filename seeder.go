package main

import (
	"log"
	"math/rand"
	"time"
)

// SeedDatabase populates the database with dummy data if it's empty
func SeedDatabase() {
	var count int64

	// ✅ Seed Users
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		users := []User{
			{Name: "User 1", Email: "user1@example.com", Password: "$2a$10$Ys8ik7V3EU.KlFZa7trJ8uSqKDsj.WGNrYu2xsAil1yT3mC.k4hwy"},
		}
		DB.Create(&users)
		log.Println("✅ Users seeded!")
	}

	// ✅ Seed Categories
	DB.Model(&Category{}).Count(&count)
	if count == 0 {
		categories := []Category{
			{Name: "Salary"},
			{Name: "Food"},
			{Name: "Transportation"},
			{Name: "Entertainment"},
			{Name: "Shopping"},
			{Name: "Investment"},
			{Name: "Bills"},
		}
		DB.Create(&categories)
		log.Println("✅ Categories seeded!")
	}

	// ✅ Seed Transactions
	DB.Model(&Transaction{}).Count(&count)
	if count == 0 {
		transactions := []Transaction{
			// Income transactions
			{Type: "Income", Amount: randomAmount(4000, 6000), Currency: "USD", ExchangeRate: 16500, Note: "Last month's salary", CategoryID: getCategoryID("Salary"), UserID: 1, CreatedAt: time.Now().AddDate(0, -6, 0)},
			{Type: "Income", Amount: randomAmount(4000, 6000), Currency: "USD", ExchangeRate: 16500, Note: "This month's salary", CategoryID: getCategoryID("Salary"), UserID: 1, CreatedAt: time.Now().AddDate(0, -5, 0)},
			{Type: "Income", Amount: randomAmount(500, 2000), Currency: "USD", ExchangeRate: 16500, Note: "Bonus", CategoryID: getCategoryID("Investment"), UserID: 1, CreatedAt: time.Now().AddDate(0, -3, 0)},

			// Expense - Food
			{Type: "Expense", Amount: randomAmount(50, 150), Currency: "USD", ExchangeRate: 16500, Note: "Lunch", CategoryID: getCategoryID("Food"), UserID: 1, CreatedAt: time.Now().AddDate(0, -6, 0)},
			{Type: "Expense", Amount: randomAmount(70, 200), Currency: "USD", ExchangeRate: 16500, Note: "Dinner outside", CategoryID: getCategoryID("Food"), UserID: 1, CreatedAt: time.Now().AddDate(0, -4, 0)},

			// Expense - Transportation
			{Type: "Expense", Amount: randomAmount(100, 300), Currency: "USD", ExchangeRate: 16500, Note: "Motorbike fuel", CategoryID: getCategoryID("Transportation"), UserID: 1, CreatedAt: time.Now().AddDate(0, -5, 0)},
			{Type: "Expense", Amount: randomAmount(50, 250), Currency: "USD", ExchangeRate: 16500, Note: "Train ticket", CategoryID: getCategoryID("Transportation"), UserID: 1, CreatedAt: time.Now().AddDate(0, -2, 0)},

			// Expense - Bills
			{Type: "Expense", Amount: randomAmount(300, 700), Currency: "USD", ExchangeRate: 16500, Note: "Electricity bill", CategoryID: getCategoryID("Bills"), UserID: 1, CreatedAt: time.Now().AddDate(0, -5, 0)},
			{Type: "Expense", Amount: randomAmount(100, 500), Currency: "USD", ExchangeRate: 16500, Note: "Monthly internet", CategoryID: getCategoryID("Bills"), UserID: 1, CreatedAt: time.Now().AddDate(0, -3, 0)},

			// Expense - Entertainment
			{Type: "Expense", Amount: randomAmount(50, 300), Currency: "USD", ExchangeRate: 16500, Note: "Movie night", CategoryID: getCategoryID("Entertainment"), UserID: 1, CreatedAt: time.Now().AddDate(0, -4, 0)},
			{Type: "Expense", Amount: randomAmount(150, 500), Currency: "USD", ExchangeRate: 16500, Note: "Online games", CategoryID: getCategoryID("Entertainment"), UserID: 1, CreatedAt: time.Now().AddDate(0, -1, 0)},
		}

		// Filter transactions with valid CategoryID
		var validTransactions []Transaction
		for _, t := range transactions {
			if t.CategoryID != nil {
				validTransactions = append(validTransactions, t)
			} else {
				log.Printf("⚠️ Transaction '%s' was not saved because CategoryID is nil!", t.Note)
			}
		}

		// Save only valid transactions
		DB.Create(&validTransactions)
		log.Println("✅ Transactions seeded!")
	}

	// ✅ Seed Budgets
	DB.Model(&Budget{}).Count(&count)
	if count == 0 {
		budgets := []Budget{
			{CategoryID: *getCategoryID("Food"), UserID: 1, Amount: 5000, Currency: "USD", ExchangeRate: 16500, Month: getCurrentMonth()},
			{CategoryID: *getCategoryID("Transportation"), UserID: 1, Amount: 3000, Currency: "USD", ExchangeRate: 16500, Month: getCurrentMonth()},
			{CategoryID: *getCategoryID("Entertainment"), UserID: 1, Amount: 2000, Currency: "USD", ExchangeRate: 16500, Month: getCurrentMonth()},
			{CategoryID: *getCategoryID("Shopping"), UserID: 1, Amount: 4000, Currency: "USD", ExchangeRate: 16500, Month: getCurrentMonth()},
			{CategoryID: *getCategoryID("Bills"), UserID: 1, Amount: 7000, Currency: "USD", ExchangeRate: 16500, Month: getCurrentMonth()},
		}

		// Filter budgets with valid CategoryID
		var validBudgets []Budget
		for _, b := range budgets {
			if b.CategoryID > 0 {
				validBudgets = append(validBudgets, b)
			} else {
				log.Printf("⚠️ Budget for category ID %v was not saved because CategoryID is nil!", b.CategoryID)
			}
		}

		// Save only valid budgets
		DB.Create(&validBudgets)
		log.Println("✅ Budgets seeded!")
	}
}

// getCategoryID retrieves the category ID by its name
func getCategoryID(name string) *uint {
	var category Category
	if err := DB.Where("name = ?", name).First(&category).Error; err == nil {
		return &category.ID
	}
	log.Printf("⚠️ Category '%s' not found!", name)
	return nil
}

// getCurrentMonth returns the current month in "YYYY-MM" format
func getCurrentMonth() string {
	return time.Now().Format("2006-01")
}

// randomAmount generates a random amount between min and max
func randomAmount(min, max int) float64 {
	return float64(rand.Intn(max-min) + min)
}
