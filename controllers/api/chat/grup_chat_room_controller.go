package chat

import (
	"net/http"
	"strconv"
	"time"

	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/chat"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateGroupRequest struct {
	Name      string   `json:"name" binding:"required"`
	MemberIDs []uint64 `json:"member_ids" binding:"required,min=1"`
}

func CreateGroupRoom(c *gin.Context) {
	var req CreateGroupRequest

	db := config.Database

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Data request tidak valid", err)
		return
	}

	if len(req.Name) < 3 {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Nama grup minimal 3 karakter", nil)
		return
	}

	userIDStr := c.GetString("user_id")
	currentUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "User tidak valid", err)
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	room := chat.ChatRoom{
		ID:        uuid.New().String(),
		Type:      "group",
		Name:      &req.Name,
		CreatedBy: currentUserID,
	}

	if err := tx.Create(&room).Error; err != nil {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat grup. Silakan hubungi administrator", err)
		return
	}

	now := time.Now()
	roomUsers := []chat.ChatRoomUser{
		{
			ID:         uuid.New().String(),
			ChatRoomID: room.ID,
			UserID:     currentUserID,
			Role:       "admin",
			JoinedAt:   &now,
		},
	}

	for _, memberID := range req.MemberIDs {
		if memberID == currentUserID {
			continue
		}

		roomUsers = append(roomUsers, chat.ChatRoomUser{
			ID:         uuid.New().String(),
			ChatRoomID: room.ID,
			UserID:     memberID,
			Role:       "member",
			JoinedAt:   &now,
		})
	}

	if err := tx.Create(&roomUsers).Error; err != nil {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan anggota ke grup. Silakan hubungi administrator", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan data grup. Silakan hubungi administrator", err)
		return
	}

	if err := db.Preload("ChatRoom").Where("chat_room_id = ?", room.ID).Find(&roomUsers).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal memuat data grup", err)
		return
	}

	helpers.SuccessResponse(c, http.StatusCreated, "Grup berhasil dibuat", gin.H{
		"room":    room,
		"members": roomUsers,
	})
}

type AddMemberRequest struct {
	ChatRoomID string   `json:"chat_room_id" binding:"required"`
	MemberIDs  []uint64 `json:"member_ids" binding:"required,min=1"`
}

func validateGroupAdmin(db *gorm.DB, roomID string, userID uint64) (*chat.ChatRoom, error) {
	var room chat.ChatRoom
	if err := db.Where("id = ? AND type = ?", roomID, "group").First(&room).Error; err != nil {
		return nil, err
	}

	var userRoom chat.ChatRoomUser
	if err := db.Where("chat_room_id = ? AND user_id = ? AND role = ?", roomID, userID, "admin").First(&userRoom).Error; err != nil {
		return nil, err
	}

	return &room, nil
}

func addNewMembers(tx *gorm.DB, roomID string, memberIDs []uint64) ([]chat.ChatRoomUser, error) {
	var addedMembers []chat.ChatRoomUser
	now := time.Now()

	for _, memberID := range memberIDs {
		var existingMember chat.ChatRoomUser
		err := tx.Where("chat_room_id = ? AND user_id = ?", roomID, memberID).First(&existingMember).Error

		if err == nil {
			continue // User sudah menjadi member
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		newMember := chat.ChatRoomUser{
			ID:         uuid.New().String(),
			ChatRoomID: roomID,
			UserID:     memberID,
			Role:       "member",
			JoinedAt:   &now,
		}

		if err := tx.Create(&newMember).Error; err != nil {
			return nil, err
		}

		addedMembers = append(addedMembers, newMember)
	}

	return addedMembers, nil
}

func AddMemberToGroup(c *gin.Context) {
	var req AddMemberRequest

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

	_, err = validateGroupAdmin(db, req.ChatRoomID, currentUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusForbidden, "Hanya admin yang dapat menambahkan anggota", err)
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal validasi admin", err)
		}
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	addedMembers, err := addNewMembers(tx, req.ChatRoomID, req.MemberIDs)
	if err != nil {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan anggota", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan data anggota. Silakan hubungi administrator", err)
		return
	}

	if len(addedMembers) == 0 {
		helpers.SuccessResponse(c, http.StatusOK, "Tidak ada anggota baru yang ditambahkan", gin.H{
			"added_members": []interface{}{},
		})
		return
	}

	var membersWithUser []chat.ChatRoomUser
	if err := db.Preload("User").Where("id IN ?", getMemberIDs(addedMembers)).Find(&membersWithUser).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal memuat data anggota", err)
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Anggota berhasil ditambahkan ke grup", gin.H{
		"room_id":       req.ChatRoomID,
		"added_members": membersWithUser,
	})
}

func getMemberIDs(members []chat.ChatRoomUser) []string {
	ids := make([]string, len(members))
	for i, member := range members {
		ids[i] = member.ID
	}
	return ids
}

type SendGroupMessageRequest struct {
	ChatRoomID string `json:"chat_room_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

func validateGroupMembership(db *gorm.DB, roomID string, userID uint64) (*chat.ChatRoom, error) {
	var room chat.ChatRoom
	if err := db.Where("id = ? AND type = ?", roomID, "group").First(&room).Error; err != nil {
		return nil, err
	}

	var userRoom chat.ChatRoomUser
	if err := db.Where("chat_room_id = ? AND user_id = ?", roomID, userID).First(&userRoom).Error; err != nil {
		return nil, err
	}

	return &room, nil
}

func createGroupMessage(db *gorm.DB, roomID string, senderID uint64, messageText string) (*chat.ChatMessage, error) {
	message := chat.ChatMessage{
		ID:         uuid.New().String(),
		ChatRoomID: roomID,
		SenderID:   senderID,
		Type:       "text",
		Message:    &messageText,
		Status:     "sent",
	}

	if err := db.Create(&message).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("ChatRoom").Preload("Sender").Where("id = ?", message.ID).First(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func SendGroupMessage(c *gin.Context) {
	var req SendGroupMessageRequest

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

	room, err := validateGroupMembership(db, req.ChatRoomID, currentUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusForbidden, "Anda bukan anggota grup ini", err)
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal validasi membership", err)
		}
		return
	}

	message, err := createGroupMessage(db, room.ID, currentUserID, req.Message)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengirim pesan", err)
		return
	}

	if err := createReceiptsForRoom(db, room.ID, message.ID, currentUserID); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mencatat status delivered", err)
		return
	}

	db.Model(&room).Update("last_message_id", message.ID)

	helpers.SuccessResponse(c, http.StatusOK, "Pesan berhasil dikirim", gin.H{
		"room_id": room.ID,
		"message": message,
	})
}
