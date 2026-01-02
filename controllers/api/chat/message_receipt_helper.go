package chat

import (
	"time"

	"mobile-directory-bussines/models/chat"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// createReceiptsForRoom mencatat status delivered untuk semua member room selain pengirim.
func createReceiptsForRoom(db *gorm.DB, roomID string, messageID string, senderID uint64) error {
	var members []chat.ChatRoomUser
	if err := db.Select("user_id").Where("chat_room_id = ? AND user_id <> ?", roomID, senderID).Find(&members).Error; err != nil {
		return err
	}

	if len(members) == 0 {
		return nil
	}

	memberIDs := make([]uint64, 0, len(members))
	for _, m := range members {
		memberIDs = append(memberIDs, m.UserID)
	}

	// Ambil receipt yang sudah ada untuk menghindari duplicate insert
	var existingUserIDs []uint64
	if err := db.Model(&chat.ChatMessageReceipt{}).
		Where("chat_message_id = ? AND user_id IN ?", messageID, memberIDs).
		Pluck("user_id", &existingUserIDs).Error; err != nil {
		return err
	}

	existing := make(map[uint64]struct{}, len(existingUserIDs))
	for _, id := range existingUserIDs {
		existing[id] = struct{}{}
	}

	now := time.Now()
	var receipts []chat.ChatMessageReceipt
	for _, uid := range memberIDs {
		if _, ok := existing[uid]; ok {
			continue
		}
		receipts = append(receipts, chat.ChatMessageReceipt{
			ID:            uuid.New().String(),
			ChatMessageID: messageID,
			UserID:        uid,
			DeliveredAt:   now,
		})
	}

	if len(receipts) == 0 {
		return nil
	}

	return db.Create(&receipts).Error
}
