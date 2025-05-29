package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
