package chat

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/chat"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMessages(c *gin.Context) {
	var chatMessages []chat.ChatMessage

	db := config.Database

	db = ApplyChatMessageFilter(c, db)

	db = helpers.GeneralSortData(db, c.Query("sort_by"), c.Query("sort_dir"), validChatMessageSortData())

	db = helpers.ApplyPreloads(db, validChatMessagePreload())

	helpers.RespondWithPagination(c, db, &chatMessages, "Gagal mengambil data chat messages")

}

func validChatMessageSortData() []string {
	return []string{
		"creted_at",
		"updated_at",
		"message",
		"type",
	}
}

func validChatMessagePreload() []string {
	return []string{
		"ChatRoom",
		"Sender",
	}
}

func ApplyChatMessageFilter(c *gin.Context, db *gorm.DB) *gorm.DB {
	if c.Query("with_deleted") == "true" {
		db = db.Unscoped()
	}

	if senderId := c.Query("sender_id"); senderId != "" {
		db = db.Where("sender_id = ?", senderId)
	}

	if chatRoomId := c.Query("chat_room_id"); chatRoomId != "" {
		db = db.Where("chat_room_id = ?", chatRoomId)
	}

	if t := c.Query("type"); t != "" {
		db = db.Where("type = ?", t)
	}

	if message := c.Query("message"); message != "" {
		db = db.Where("message LIKE ?", "%"+message+"%")
	}

	return db
}

type UpdateMessageRequest struct {
	Message *string `json:"message"`
	Type    string  `json:"type"`
}

func UpdateMessage(c *gin.Context) {
	id := c.Param("id")

	var req UpdateMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, 400, "Data request tidak valid", err)
		return
	}

	var message chat.ChatMessage

	if err := config.Database.First(&message, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, 404, "Pesan tidak ditemukan", nil)
		} else {
			helpers.ErrorResponse(c, 500, "Gagal mengambil data pesan", err)
		}
		return
	}

	if req.Message != nil {
		message.Message = req.Message
	}
	if req.Type != "" {
		message.Type = req.Type
	}

	if err := config.Database.Save(&message).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal update pesan", err)
		return
	}

	if err := config.Database.Preload("ChatRoom").Preload("Sender").First(&message, "id = ?", id).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal memuat data pesan", err)
		return
	}

	helpers.SuccessResponse(c, 200, "Pesan berhasil diupdate", message)
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")

	var message chat.ChatMessage
	if err := config.Database.First(&message, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, 404, "Pesan tidak ditemukan", nil)
		} else {
			helpers.ErrorResponse(c, 500, "Gagal mengambil data pesan", err)
		}
		return
	}

	if err := config.Database.Delete(&message).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal menghapus pesan", err)
		return
	}

	helpers.SuccessResponse(c, 200, "Pesan berhasil dihapus", nil)
}

func DeleteChatMessage(c *gin.Context) {
	id := c.Param("id")

	var chatMessage chat.ChatMessage

	db := config.Database

	if err := db.First(&chatMessage, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, 404, "Pesan tidak ditemukan", nil)
		} else {
			helpers.ErrorResponse(c, 500, "Gagal mengambil data pesan", err)
		}

		return
	}

	if err := db.Delete(&chatMessage).Error; err != nil {
		helpers.ErrorResponse(c, 500, "Gagal menghapus pesan", err)
		return
	}

	helpers.SuccessResponse(c, 200, "Pesan berhasil dihapus", 204)
}
