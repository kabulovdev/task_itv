package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `json:"title" binding:"required"`
	Director string `json:"director" binding:"required"`
	Year     int    `json:"year" binding:"required"`
	Plot     string `json:"plot"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Movie{})
}
