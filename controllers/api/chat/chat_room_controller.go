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

func GetChatRooms(c *gin.Context) {
	var chatRooms []chat.ChatRoom

	db := config.Database

	db = ApplyChatRoomFilter(c, db)

	db = helpers.GeneralSortData(db, c.Query("sort_by"), c.Query("sort_dir"), validSortData())

	db = helpers.ApplyPreloads(db, validPreload())

	helpers.RespondWithPagination(c, db, &chatRooms, "Gagal mengambil data chat rooms")
}

func validSortData() []string {
	return []string{
		"created_at",
		"updated_at",
		"name",
		"type",
	}
}

func validPreload() []string {
	return []string{}
}

func ApplyChatRoomFilter(c *gin.Context, db *gorm.DB) *gorm.DB {
	if t := c.Query("type"); t != "" {
		db = db.Where("type = ?", t)
	}

	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if startCreatedAt := c.Query("created_at"); startCreatedAt != "" {
		db = db.Where("created_at >= ?", startCreatedAt)
	}

	if finishCreatedAt := c.Query("created_at"); finishCreatedAt != "" {
		db = db.Where("created_at <= ?", finishCreatedAt)
	}

	return db
}

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

	if req.TargetUserID == currentUserID {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Tidak dapat mengirim pesan ke diri sendiri", nil)
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
		Status:     "sent",
	}
	if err := db.Create(&message).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengirim pesan. Silakan hubungi administrator", err)
		return
	}

	// Catat delivered untuk semua member room selain pengirim
	if err := createReceiptsForRoom(db, room.ID, message.ID, currentUserID); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mencatat status delivered", err)
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
