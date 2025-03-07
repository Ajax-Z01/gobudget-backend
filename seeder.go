package main

import "log"

func SeedDatabase() {
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		users := []User{
			{Name: "User 1", Email: "user1@example.com", Password: "password123"},
		}
		DB.Create(&users)
		log.Println("✅ Users seeded!")
	}
	DB.Model(&Category{}).Count(&count)
	if count == 0 {
		categories := []Category{
			{Name: "Gaji"},
			{Name: "Makanan"},
			{Name: "Transportasi"},
			{Name: "Hiburan"},
		}
		DB.Create(&categories)
		log.Println("✅ Categories seeded!")
	}

	DB.Model(&Transaction{}).Count(&count)
	if count == 0 {
		transactions := []Transaction{
			{Type: "Income", Amount: 5000, Note: "Gaji bulan ini", CategoryID: getCategoryID("Gaji"), UserID: 1},
			{Type: "Expense", Amount: 100, Note: "Makan siang", CategoryID: getCategoryID("Makanan"), UserID: 1},
			{Type: "Expense", Amount: 200, Note: "Beli bensin", CategoryID: getCategoryID("Transportasi"), UserID: 1},
		}
		DB.Create(&transactions)
		log.Println("✅ Transactions seeded!")
	}
}

func getCategoryID(name string) *uint {
	var category Category
	if err := DB.Where("name = ?", name).First(&category).Error; err == nil {
		return &category.ID
	}
	return nil
}
