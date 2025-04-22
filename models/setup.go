package models

import (
	"fmt"// this helps to use print statements
	"log"// llog is used to log errors
	"os" // This is the missing import
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	os.Getenv("DB_USERNAME"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_DATABASE"),
)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	database.AutoMigrate(&Note{})
	DB = database
	fmt.Println("Database connected")
}
