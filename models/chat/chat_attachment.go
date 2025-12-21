package chat

import (
	"time"
)

type ChatAttachment struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	ChatMessageID string    `gorm:"not null" json:"chat_message_id"`
	FileName      string    `gorm:"type:varchar(255)" json:"file_name"`
	FileURL       string    `gorm:"type:varchar(500)" json:"file_url"`
	FileType      string    `gorm:"type:varchar(50)" json:"file_type"`
	FileSize      int64     `json:"file_size"`
	CreatedAt     time.Time `json:"created_at"`

	ChatMessage ChatMessage `gorm:"foreignKey:ChatMessageID" json:"chat_message,omitempty"`
}

func (ChatAttachment) TableName() string {
	return "chat_attachments"
}
