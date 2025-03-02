package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Model To-Do
type Todo struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

// Koneksi ke database
func InitDB() {
	dsn := "host=localhost user=postgres password=1234 dbname=todo_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	DB = db
	DB.AutoMigrate(&Todo{})
}

// Handler untuk mendapatkan semua to-do
func GetTodos(c *gin.Context) {
	var todos []Todo
	DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

// Handler untuk menambahkan to-do
func CreateTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&todo)
	c.JSON(http.StatusCreated, todo)
}

// Handler untuk memperbarui to-do
func UpdateTodo(c *gin.Context) {
	var todo Todo
	id := c.Param("id")
	if err := DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Save(&todo)
	c.JSON(http.StatusOK, todo)
}

// Handler untuk menghapus to-do
func DeleteTodo(c *gin.Context) {
	var todo Todo
	id := c.Param("id")
	if err := DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	DB.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func main() {
	InitDB()

	r := gin.Default()

	r.GET("/todos", GetTodos)
	r.POST("/todos", CreateTodo)
	r.PUT("/todos/:id", UpdateTodo)
	r.DELETE("/todos/:id", DeleteTodo)

	r.Run(":8080") // Jalankan server di port 8080
}
