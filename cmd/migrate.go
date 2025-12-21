package main

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/models"
	"mobile-directory-bussines/models/chat"
)

func Migrate() {
	config.Database.AutoMigrate(
		&models.User{},
		&chat.ChatRoom{},
		&chat.ChatRoomUser{},
		&chat.ChatMessage{},
		&chat.ChatMessageRead{},
		&chat.ChatMessageDelete{},
		&chat.ChatAttachment{},
	)
}

func main() {
	config.ConnectDatabase()
	Migrate()
}
