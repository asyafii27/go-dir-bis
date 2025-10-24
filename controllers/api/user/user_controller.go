package controllers

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models"
	"mobile-directory-bussines/models/role"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	db := helpers.PreloadIfParam(c, config.Database, "role", "Role")

	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if roleID := c.Query("role_id"); roleID != "" {
		db = db.Where("role_id = ?", roleID)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	allowedSorts := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"role_id":    true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	db = db.Order(sortBy + " " + sortDir)

	meta, err := helpers.Paginate(c, db, &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if c.Query("role") == "" {
		for i := range users {
			users[i].Role = nil
		}
	}

	response := gin.H{
		"status": "success",
		"data":   users,
	}

	if meta.Page > 0 {
		response["meta"] = meta
	}

	c.JSON(http.StatusOK, response)
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := config.Database.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

func StoreUser(c *gin.Context) {
	var input struct {
		Name                 string `json:"name" binding:"required,max=255"`
		Email                string `json:"email" binding:"required,email,max=255"`
		Password             string `json:"password" binding:"required,min=8"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required,min=8"`
		RoleID               uint   `json:"role_id" binding:"required"`
	}

	if ok := validateCreateUserInput(c, &input); !ok {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   input.RoleID,
	}

	if err := config.Database.Create(&user).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat user")
		return
	}

	helpers.SuccessResponse(c, http.StatusCreated, "Success", user)
}

func validateCreateUserInput(c *gin.Context, input *struct {
	Name                 string `json:"name" binding:"required,max=255"`
	Email                string `json:"email" binding:"required,email,max=255"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,min=8"`
	RoleID               uint   `json:"role_id" binding:"required"`
}) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		validationErrorResponse(c, err)
		return false
	}

	if input.Password != input.PasswordConfirmation || !validPassword(input.Password) {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Password harus minimal 8 karakter, mengandung huruf besar, kecil, angka, dan sesuai konfirmasi")
		return false
	}

	var existingUser models.User
	if err := config.Database.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Email sudah digunakan")
		return false
	}

	var role role.Role
	if err := config.Database.First(&role, input.RoleID).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Role tidak ditemukan")
		return false
	}

	return true
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.Database.First(&user, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan")
		return
	}

	var input struct {
		Name                 string `json:"name" binding:"required,max=255"`
		Email                string `json:"email" binding:"required,email,max=255"`
		Password             string `json:"password"` // opsional
		PasswordConfirmation string `json:"password_confirmation"`
		RoleID               uint   `json:"role_id" binding:"required"`
	}

	if ok := validateUpdateUserInput(c, &input, user.ID); !ok {
		return
	}

	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	user.Name = input.Name
	user.Email = input.Email
	user.RoleID = input.RoleID

	if err := config.Database.Save(&user).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate user")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "User berhasil diperbarui", user)
}

func validateUpdateUserInput(c *gin.Context, input *struct {
	Name                 string `json:"name" binding:"required,max=255"`
	Email                string `json:"email" binding:"required,email,max=255"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	RoleID               uint   `json:"role_id" binding:"required"`
}, userID uint) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		validationErrorResponse(c, err)
		return false
	}

	if input.Password != "" {
		if input.Password != input.PasswordConfirmation || !validPassword(input.Password) {
			helpers.ErrorResponse(c, http.StatusBadRequest, "Password harus minimal 8 karakter, mengandung huruf besar, kecil, angka, dan sesuai konfirmasi")
			return false
		}
	}

	var existingUser models.User
	if err := config.Database.Where("email = ? AND id != ?", input.Email, userID).First(&existingUser).Error; err == nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Email sudah digunakan oleh user lain")
		return false
	}

	var role role.Role
	if err := config.Database.First(&role, input.RoleID).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Role tidak ditemukan")
		return false
	}

	return true
}

func validPassword(pw string) bool {
	lower := regexp.MustCompile(`[a-z]`).MatchString
	upper := regexp.MustCompile(`[A-Z]`).MatchString
	num := regexp.MustCompile(`\d`).MatchString
	return len(pw) >= 8 && lower(pw) && upper(pw) && num(pw)
}

func validationErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": "Validasi gagal",
		"error":   err.Error(),
	})
}
