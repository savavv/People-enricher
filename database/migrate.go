package database

import (
	"log"
	"people-enricher/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
