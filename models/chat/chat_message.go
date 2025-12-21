package chat

import (
	"time"

	"mobile-directory-bussines/models"

	"gorm.io/gorm"
)

type ChatMessage struct {
	ID         string         `gorm:"type:char(36);primaryKey" json:"id"`
	ChatRoomID string         `gorm:"not null" json:"chat_room_id"`
	SenderID   uint64         `gorm:"not null" json:"sender_id"`
	Type       string         `gorm:"type:enum('text','image','file','system');default:'text'" json:"type"`
	Message    *string        `gorm:"type:text" json:"message"`
	FileURL    *string        `gorm:"type:varchar(500)" json:"file_url"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	ChatRoom ChatRoom    `gorm:"foreignKey:ChatRoomID" json:"chat_room,omitempty"`
	Sender   models.User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
