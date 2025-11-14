package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDatabase() {
	var databaseUser = "root"
	var databasePassword = ""
	var databaseHost = "127.0.0.1"
	var databasePort = "3306"
	var databaseName = "directory_business"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database: ", err)
	}

	fmt.Println("Database Connected!")
	Database = db
}
