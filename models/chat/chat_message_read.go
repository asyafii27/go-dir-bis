package chat

import (
	"time"

	"mobile-directory-bussines/models"
)

type ChatMessageRead struct {
	ID            string     `gorm:"type:char(36);primaryKey" json:"id"`
	ChatMessageID string     `gorm:"not null;uniqueIndex:uniq_message_user" json:"chat_message_id"`
	UserID        uint64     `gorm:"not null;uniqueIndex:uniq_message_user" json:"user_id"`
	ReadAt        *time.Time `json:"read_at"`

	ChatMessage ChatMessage `gorm:"foreignKey:ChatMessageID" json:"chat_message,omitempty"`
	User        models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ChatMessageRead) TableName() string {
	return "chat_message_reads"
}
