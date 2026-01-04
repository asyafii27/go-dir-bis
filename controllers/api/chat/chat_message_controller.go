package chat

import (
	"fmt"
	"strconv"
	"time"

	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/chat"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMessages(c *gin.Context) {
	var chatMessages []chat.ChatMessage

	db := config.Database

	// Ambil user_id dari context (user yang sedang membaca pesan)
	userIDStr := c.GetString("user_id")
	currentUserID, _ := strconv.ParseUint(userIDStr, 10, 64)
	chatRoomID := c.Query("chat_room_id")

	// Mark messages as read SEBELUM query (agar data yang diambil sudah terupdate)
	if chatRoomID != "" && currentUserID > 0 {
		markMessagesAsRead(config.Database, chatRoomID, currentUserID)
	}

	db = ApplyChatMessageFilter(c, db)

	db = helpers.GeneralSortData(db, c.Query("sort_by"), c.Query("sort_dir"), validChatMessageSortData())

	db = helpers.ApplyPreloads(db, validChatMessagePreload())

	helpers.RespondWithPagination(c, db, &chatMessages, "Gagal mengambil data chat messages")

	// Debug: Lihat messages yang diambil
	fmt.Println("\n=== DEBUG GET MESSAGES ===")
	fmt.Printf("Current User ID: %d\n", currentUserID)
	fmt.Printf("Chat Room ID: %s\n", chatRoomID)
	fmt.Printf("Total Messages: %d\n", len(chatMessages))
	for i, msg := range chatMessages {
		fmt.Printf("[%d] ID: %s, Sender: %d, Status: %s\n", i, msg.ID, msg.SenderID, msg.Status)
	}
	fmt.Println("=== END DEBUG ===")
}

// markMessagesAsRead marks all messages in a room as read for a specific user
func markMessagesAsRead(db *gorm.DB, roomID string, userID uint64) {
	fmt.Println("\n=== DEBUG MARK AS READ ===")
	fmt.Printf("Room ID: %s, Reader User ID: %d\n", roomID, userID)

	// Ambil semua message ID di room ini yang bukan dikirim oleh user ini
	var messageIDs []string
	db.Model(&chat.ChatMessage{}).
		Select("id").
		Where("chat_room_id = ? AND sender_id != ? AND status != ?", roomID, userID, "read").
		Pluck("id", &messageIDs)

	fmt.Printf("Messages to mark as read: %d\n", len(messageIDs))
	spew.Dump(messageIDs)

	if len(messageIDs) == 0 {
		fmt.Println("No messages to mark")
		fmt.Println("=== END DEBUG ===")
		return
	}

	// Cek message mana yang belum ada read receipt untuk user ini
	var existingReadMessageIDs []string
	db.Model(&chat.ChatMessageRead{}).
		Select("chat_message_id").
		Where("user_id = ? AND chat_message_id IN ?", userID, messageIDs).
		Pluck("chat_message_id", &existingReadMessageIDs)

	// Buat map untuk cepat cek existing
	existingMap := make(map[string]bool)
	for _, id := range existingReadMessageIDs {
		existingMap[id] = true
	}

	// Create read receipts untuk message yang belum dibaca
	now := time.Now()
	var readReceipts []chat.ChatMessageRead
	var unreadMessageIDs []string

	for _, msgID := range messageIDs {
		if !existingMap[msgID] {
			readReceipts = append(readReceipts, chat.ChatMessageRead{
				ID:            uuid.New().String(),
				ChatMessageID: msgID,
				UserID:        userID,
				ReadAt:        &now,
			})
			unreadMessageIDs = append(unreadMessageIDs, msgID)
		}
	}

	// Gunakan transaction untuk memastikan atomicity
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if len(readReceipts) > 0 {
		fmt.Printf("Creating %d read receipts\n", len(readReceipts))

		// Create read receipts
		if err := tx.Create(&readReceipts).Error; err != nil {
			fmt.Printf("Error creating read receipts: %v\n", err)
			tx.Rollback()
			return
		}

		fmt.Printf("Updating %d messages to 'read' status\n", len(unreadMessageIDs))

		// Update status message menjadi "read"
		if err := tx.Model(&chat.ChatMessage{}).
			Where("id IN ?", unreadMessageIDs).
			Update("status", "read").Error; err != nil {
			fmt.Printf("Error updating status: %v\n", err)
			tx.Rollback()
			return
		}

		fmt.Println("✓ Successfully marked messages as read")
	}

	tx.Commit()
	fmt.Println("=== END DEBUG ===")
}

func validChatMessageSortData() []string {
	return []string{
		"creted_at",
		"updated_at",
		"message",
		"type",
	}
}

func validChatMessagePreload() []string {
	return []string{
		"ChatRoom",
		"Sender",
	}
}

func ApplyChatMessageFilter(c *gin.Context, db *gorm.DB) *gorm.DB {
	if c.Query("with_deleted") == "true" {
		db = db.Unscoped()
	}

	if senderId := c.Query("sender_id"); senderId != "" {
		db = db.Where("sender_id = ?", senderId)
	}

	if chatRoomId := c.Query("chat_room_id"); chatRoomId != "" {
		db = db.Where("chat_room_id = ?", chatRoomId)
	}

	if t := c.Query("type"); t != "" {
		db = db.Where("type = ?", t)
	}

	if message := c.Query("message"); message != "" {
		db = db.Where("message LIKE ?", "%"+message+"%")
	}

	return db
}

type UpdateMessageRequest struct {
	Message *string `json:"message"`
	Type    string  `json:"type"`
}

func UpdateMessage(c *gin.Context) {
	id := c.Param("id")

	var req UpdateMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, 400, "Data request tidak valid", err)
		return
	}

	var message chat.ChatMessage

	if err := config.Database.First(&message, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, 404, "Pesan tidak ditemukan", nil)
		} else {
			helpers.ErrorResponse(c, 500, "Gagal mengambil data pesan", err)
		}
		return
	}

	if req.Message != nil {
		message.Message = req.Message
	}
	if req.Type != "" {
		message.Type = req.Type
	}

	if err := config.Database.Save(&message).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal update pesan", err)
		return
	}

	if err := config.Database.Preload("ChatRoom").Preload("Sender").First(&message, "id = ?", id).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal memuat data pesan", err)
		return
	}

	helpers.SuccessResponse(c, 200, "Pesan berhasil diupdate", message)
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")

	var message chat.ChatMessage
	if err := config.Database.First(&message, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, 404, "Pesan tidak ditemukan", nil)
		} else {
			helpers.ErrorResponse(c, 500, "Gagal mengambil data pesan", err)
		}
		return
	}

	if err := config.Database.Delete(&message).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal menghapus pesan", err)
		return
	}

	helpers.SuccessResponse(c, 200, "Pesan berhasil dihapus", nil)
}

func DeleteChatMessage(c *gin.Context) {
	id := c.Param("id")

	var chatMessage chat.ChatMessage

	db := config.Database

	if err := db.First(&chatMessage, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, 404, "Pesan tidak ditemukan", nil)
		} else {
			helpers.ErrorResponse(c, 500, "Gagal mengambil data pesan", err)
		}

		return
	}

	if err := db.Delete(&chatMessage).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal menghapus pesan", err)
		return
	}

	helpers.SuccessResponse(c, 200, "Pesan berhasil dihapus", 204)
}
