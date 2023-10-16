package database

import (
	"fiber-project/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {

	// connect to database
	db := sqlite.Open("bin/database/users.db")
	gormDB, err := gorm.Open(db, &gorm.Config{})

	if err != nil {
		panic("Error connecting to gorm")
	}
	DB = gormDB
	gormDB.AutoMigrate(&models.User{})
}
