package chat

import (
	"time"

	"mobile-directory-bussines/models"
)

type ChatMessageDelete struct {
	ID            string     `gorm:"type:char(36);primaryKey" json:"id"`
	ChatMessageID string     `gorm:"not null;uniqueIndex:uniq_message_user" json:"chat_message_id"`
	UserID        uint64     `gorm:"not null;uniqueIndex:uniq_message_user" json:"user_id"`
	DeletedAt     *time.Time `json:"deleted_at"`

	ChatMessage ChatMessage `gorm:"foreignKey:ChatMessageID" json:"chat_message,omitempty"`
	User        models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ChatMessageDelete) TableName() string {
	return "chat_message_deletes"
}
