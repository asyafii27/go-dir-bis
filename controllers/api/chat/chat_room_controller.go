package chat

import (
	"net/http"
	"strconv"

	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/chat"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateMessageRequest struct {
	TargetUserID uint64 `json:"target_user_id" binding:"required"`
	Message      string `json:"message" binding:"required"`
}

func StorePrivateMessage(c *gin.Context) {
	var req PrivateMessageRequest

	db := config.Database

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Data request tidak valid", err)
		return
	}

	userIDStr := c.GetString("user_id")
	currentUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "User tidak valid", err)
		return
	}

	var existingRoom chat.ChatRoom
	subQuery := db.Model(&chat.ChatRoomUser{}).
		Select("chat_room_id").
		Where("user_id IN (?, ?)", currentUserID, req.TargetUserID).
		Group("chat_room_id").
		Having("COUNT(DISTINCT user_id) = 2")

	err = db.Where("type = ? AND id IN (?)", "private", subQuery).
		First(&existingRoom).Error

	var room chat.ChatRoom
	if err == gorm.ErrRecordNotFound {
		room = chat.ChatRoom{
			ID:        uuid.New().String(),
			Type:      "private",
			CreatedBy: currentUserID,
		}
		if err := db.Create(&room).Error; err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat room. Silakan hubungi administrator", err)
			return
		}

		roomUsers := []chat.ChatRoomUser{
			{ID: uuid.New().String(), ChatRoomID: room.ID, UserID: currentUserID, Role: "member"},
			{ID: uuid.New().String(), ChatRoomID: room.ID, UserID: req.TargetUserID, Role: "member"},
		}
		if err := db.Create(&roomUsers).Error; err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan user ke room. Silakan hubungi administrator", err)
			return
		}
	} else if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengirim pesan. Silakan hubungi administrator", err)
		return
	} else {
		room = existingRoom
	}

	message := chat.ChatMessage{
		ID:         uuid.New().String(),
		ChatRoomID: room.ID,
		SenderID:   currentUserID,
		Type:       "text",
		Message:    &req.Message,
	}
	if err := db.Create(&message).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengirim pesan. Silakan hubungi administrator", err)
		return
	}

	if err := db.Preload("ChatRoom").
		Preload("Sender").
		Where("id = ?", message.ID).
		First(&message).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal memuat data pesan", err)
		return
	}

	if err := db.Model(&room).Update("last_message_id", message.ID).Error; err != nil {
	}

	helpers.SuccessResponse(c, http.StatusOK, "Message sent successfully", gin.H{
		"room_id": room.ID,
		"message": message,
	})
}
