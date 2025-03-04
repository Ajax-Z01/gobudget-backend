package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitTestDatabase() {
	dsn := "host=localhost user=postgres password=24'K>W6tMr5an!? dbname=test_db port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	DB.AutoMigrate(&User{}, &Category{}, &Transaction{})
	fmt.Println("Test database connected & migrated successfully!")
}

func SetupTestRouter() *gin.Engine {
	InitTestDatabase()
	CleanupDatabase()
	SeedDatabase()

	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/transactions", GetTransactions)
	r.GET("/categories", GetCategories)

	return r
}

func CleanupDatabase() {
	tables := []string{"users"}
	for _, table := range tables {
		err := DB.Exec(fmt.Sprintf(`DELETE FROM %s`, table)).Error
		if err != nil {
			log.Printf("Failed to clean up table %s: %v", table, err)
		}
	}
}

func TestRegister(t *testing.T) {
	r := SetupTestRouter()

	body := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	r := SetupTestRouter()

	body := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	reqReg, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(body)))
	reqReg.Header.Set("Content-Type", "application/json")

	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	assert.Equal(t, http.StatusCreated, wReg.Code)

	body = `{"email":"test@example.com","password":"password123"}`
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	reqLogin.Header.Set("Content-Type", "application/json")

	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusOK, wLogin.Code)
}

func TestGetTransactions(t *testing.T) {
	r := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/transactions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, "[]", w.Body.String())
}

func TestGetCategories(t *testing.T) {
	r := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, "[]", w.Body.String())
}
