package chat

import (
	"time"

	"mobile-directory-bussines/models"
)

// ChatMessageReceipt merekam pesan yang sudah diterima (delivered) oleh user, sebelum dibaca.
type ChatMessageReceipt struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	ChatMessageID string    `gorm:"not null;uniqueIndex:uniq_receipt_message_user" json:"chat_message_id"`
	UserID        uint64    `gorm:"not null;uniqueIndex:uniq_receipt_message_user" json:"user_id"`
	DeliveredAt   time.Time `json:"delivered_at"`

	ChatMessage ChatMessage `gorm:"foreignKey:ChatMessageID" json:"chat_message,omitempty"`
	User        models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ChatMessageReceipt) TableName() string {
	return "chat_message_receipts"
}
