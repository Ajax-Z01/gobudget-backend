package main

import "log"

func SeedDatabase() {
	var count int64
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
			{Type: "Income", Amount: 5000, Note: "Gaji bulan ini", CategoryID: getCategoryID("Gaji")},
			{Type: "Expense", Amount: 1000, Note: "Makan siang", CategoryID: getCategoryID("Makanan")},
			{Type: "Expense", Amount: 2000, Note: "Bensin motor", CategoryID: getCategoryID("Transportasi")},
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
