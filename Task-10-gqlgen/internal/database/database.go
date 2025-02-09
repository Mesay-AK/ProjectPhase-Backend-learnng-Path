package database

import (
	"fmt"
	"log"
	// "os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"taskManager/internal/models"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=admin password=admin dbname=taskmanager port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB = db
	fmt.Println("Database connected successfully!")
}

func Migrate() {
	err := DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Database migrated successfully!")
}
