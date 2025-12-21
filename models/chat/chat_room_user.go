package chat

import (
	"mobile-directory-bussines/models"
	"time"

	"gorm.io/gorm"
)

type ChatRoomUser struct {
	ID                string         `gorm:"type:char(36);primaryKey" json:"id"`
	ChatRoomID        string         `gorm:"not null;uniqueIndex:uniq_room_user" json:"chat_room_id"`
	UserID            uint64         `gorm:"not null;uniqueIndex:uniq_room_user" json:"user_id"`
	Role              string         `gorm:"type:enum('member','admin');default:'member'" json:"role"`
	LastReadMessageID *uint64        `json:"last_read_message_id"`
	IsMuted           bool           `gorm:"default:false" json:"is_muted"`
	JoinedAt          *time.Time     `json:"joined_at"`
	LeftAt            *time.Time     `json:"left_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	ChatRoom ChatRoom    `gorm:"foreignKey:ChatRoomID" json:"chat_room,omitempty"`
	User     models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ChatRoomUser) TableName() string {
	return "chat_room_users"
}
