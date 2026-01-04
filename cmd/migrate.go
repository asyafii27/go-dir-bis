package main

import (
	"fmt"
	"log"

	"mobile-directory-bussines/config"
	"mobile-directory-bussines/migrations"
	"mobile-directory-bussines/models"
	"mobile-directory-bussines/models/chat"
)

func Migrate() {
	config.Database.AutoMigrate(
		&models.User{},
		&chat.ChatRoom{},
		&chat.ChatRoomUser{},
		&chat.ChatMessage{},
		&chat.ChatMessageReceipt{},
		&chat.ChatMessageRead{},
		&chat.ChatMessageDelete{},
		&chat.ChatAttachment{},
	)
}

func RunCustomMigrations() {
	db := config.Database

	// Migration: Update enum status untuk menambahkan 'read'
	if err := migrations.AddReadStatusToMessages(db); err != nil {
		log.Printf("Warning: %v\n", err)
	}
}

func main() {
	config.ConnectDatabase()

	fmt.Println("=== Running Migrations ===")
	Migrate()
	fmt.Println("✓ AutoMigrate selesai")

	fmt.Println("\n=== Running Custom Migrations ===")
	RunCustomMigrations()

	fmt.Println("\n=== Migration Complete ===")
}
