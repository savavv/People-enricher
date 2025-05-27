package database

import (
	"log"
	"os"
	"people-enricher/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("POSTGRES_DSN")
	log.Printf("[INFO] Connecting to DB with DSN: %s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to database: %v", err)
	}
	DB = db

	log.Println("[INFO] Running migrations...")
	err = DB.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatalf("[ERROR] failed to migrate database: %v", err)
	}
	log.Println("[INFO] Database migration completed")
}
