package chat

import (
	"fmt"
	"time"

	"mobile-directory-bussines/models/chat"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// markMessagesAsDelivered marks messages as delivered when user opens room
func markMessagesAsDelivered(db *gorm.DB, roomID string, userID uint64) error {
	fmt.Println("\n=== DEBUG MARK AS DELIVERED ===")
	fmt.Printf("Room ID: %s, User ID: %d\n", roomID, userID)

	// Get room info untuk cek tipe
	var room chat.ChatRoom
	if err := db.First(&room, "id = ?", roomID).Error; err != nil {
		return fmt.Errorf("gagal get room: %v", err)
	}

	// Ambil semua message di room ini yang bukan dikirim oleh user ini
	var messages []chat.ChatMessage
	db.Where("chat_room_id = ? AND sender_id != ?", roomID, userID).
		Find(&messages)

	if len(messages) == 0 {
		fmt.Println("No messages to mark as delivered")
		fmt.Println("=== END DEBUG ===")
		return nil
	}

	fmt.Printf("Found %d messages to process\n", len(messages))

	// Collect message IDs untuk create receipts
	var messageIDs []string
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}

	// Cek message mana yang belum ada receipt untuk user ini
	var existingReceiptMessageIDs []string
	db.Model(&chat.ChatMessageReceipt{}).
		Select("chat_message_id").
		Where("user_id = ? AND chat_message_id IN ?", userID, messageIDs).
		Pluck("chat_message_id", &existingReceiptMessageIDs)

	// Buat map untuk cepat cek existing
	existingMap := make(map[string]bool)
	for _, id := range existingReceiptMessageIDs {
		existingMap[id] = true
	}

	// Create receipts untuk message yang belum ada
	now := time.Now()
	var receipts []chat.ChatMessageReceipt
	var newReceiptMessageIDs []string

	for _, msgID := range messageIDs {
		if !existingMap[msgID] {
			receipts = append(receipts, chat.ChatMessageReceipt{
				ID:            uuid.New().String(),
				ChatMessageID: msgID,
				UserID:        userID,
				DeliveredAt:   now,
			})
			newReceiptMessageIDs = append(newReceiptMessageIDs, msgID)
		}
	}

	if len(receipts) == 0 {
		fmt.Println("All messages already have receipts")
		fmt.Println("=== END DEBUG ===")
		return nil
	}

	fmt.Printf("Creating %d new receipts\n", len(receipts))

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create receipts
	if err := tx.Create(&receipts).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal create receipts: %v", err)
	}

	// Update status message berdasarkan tipe room
	if room.Type == "private" {
		// Untuk private chat, langsung update ke delivered
		fmt.Printf("Private chat: updating %d messages to delivered\n", len(newReceiptMessageIDs))
		if err := tx.Model(&chat.ChatMessage{}).
			Where("id IN ? AND status = ?", newReceiptMessageIDs, "sent").
			Update("status", "delivered").Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update status: %v", err)
		}
	} else if room.Type == "group" {
		// Untuk group chat, cek apakah semua member sudah terima
		fmt.Println("Group chat: checking if all members received")

		// Get total members (exclude sender)
		for _, msgID := range newReceiptMessageIDs {
			var msg chat.ChatMessage
			if err := tx.First(&msg, "id = ?", msgID).Error; err != nil {
				continue
			}

			var totalMembers int64
			tx.Model(&chat.ChatRoomUser{}).
				Where("chat_room_id = ? AND user_id != ?", roomID, msg.SenderID).
				Count(&totalMembers)

			var totalReceipts int64
			tx.Model(&chat.ChatMessageReceipt{}).
				Where("chat_message_id = ?", msgID).
				Count(&totalReceipts)

			fmt.Printf("Message %s: %d/%d members received\n", msgID, totalReceipts, totalMembers)

			// Jika semua member sudah terima, update status ke delivered
			if totalReceipts >= totalMembers && msg.Status == "sent" {
				fmt.Printf("All members received message %s, updating to delivered\n", msgID)
				if err := tx.Model(&chat.ChatMessage{}).
					Where("id = ?", msgID).
					Update("status", "delivered").Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("gagal update status grup: %v", err)
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit: %v", err)
	}

	fmt.Println("✓ Successfully marked messages as delivered")
	fmt.Println("=== END DEBUG ===")
	return nil
}

// updateGroupMessageReadStatus updates message status to read if all members have read it
func updateGroupMessageReadStatus(db *gorm.DB, messageID string) error {
	var msg chat.ChatMessage
	if err := db.First(&msg, "id = ?", messageID).Error; err != nil {
		return err
	}

	// Get total members (exclude sender)
	var totalMembers int64
	db.Model(&chat.ChatRoomUser{}).
		Where("chat_room_id = ? AND user_id != ?", msg.ChatRoomID, msg.SenderID).
		Count(&totalMembers)

	// Get total reads
	var totalReads int64
	db.Model(&chat.ChatMessageRead{}).
		Where("chat_message_id = ?", messageID).
		Count(&totalReads)

	fmt.Printf("Group message %s: %d/%d members read\n", messageID, totalReads, totalMembers)

	// Jika semua member sudah baca, update status ke read
	if totalReads >= totalMembers && msg.Status != "read" {
		fmt.Printf("All members read message %s, updating to read\n", messageID)
		if err := db.Model(&chat.ChatMessage{}).
			Where("id = ?", messageID).
			Update("status", "read").Error; err != nil {
			return err
		}
	}

	return nil
}
