package chat

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {
	ID            string         `gorm:"type:char(36);primaryKey" json:"id"`
	Type          string         `gorm:"type:enum('private','group');not null" json:"type"`
	Name          *string        `gorm:"type:varchar(255)" json:"name"`
	CreatedBy     uint64         `gorm:"not null" json:"created_by"`
	LastMessageID *uint64        `json:"last_message_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
