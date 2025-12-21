package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var databaseUser = os.Getenv("DB_USER")
	if databaseUser == "" {
		databaseUser = "root"
	}
	var databasePassword = os.Getenv("DB_PASSWORD")
	var databaseHost = os.Getenv("DB_HOST")
	if databaseHost == "" {
		databaseHost = "127.0.0.1"
	}
	var databasePort = os.Getenv("DB_PORT")
	if databasePort == "" {
		databasePort = "3306"
	}
	var databaseName = os.Getenv("DB_NAME")
	if databaseName == "" {
		databaseName = "directory_business_go"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database: ", err)
	}

	fmt.Println("Database Connected!")
	Database = db
}
