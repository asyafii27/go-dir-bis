package category

import (
	"fmt"
	"log"
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/master/category"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context) {
	var categories []category.Category
	db := config.Database

	db = ApplyCategoryFilters(c, db)

	page := c.Query("page")

	if page != "" {
		meta, err := helpers.LaravelPaginate(c, db, &categories)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, meta)
		return
	}

	if err := db.Find(&categories).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, categories)
}

func ApplyCategoryFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ", "%"+name+"%")
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":   true,
		"name": true,
	}

	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return db.Order(sortBy + " " + sortDir)
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category category.Category

	if err := config.Database.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "Category tidak ditemukan")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", category)
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"required,max=50"`
}

func generateCategoryCode() string {
	var lastCategory category.Category
	result := config.Database.Order("id desc").First(&lastCategory)

	nextID := uint64(1)
	if result.Error == nil {
		nextID = lastCategory.ID + 1
	}

	return fmt.Sprintf("CAT%06d", nextID)
}

func StoreCategory(c *gin.Context) {
	var request CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationResponse := helpers.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationResponse)
		return
	}

	db := config.Database

	tx := db.Begin()
	if tx.Error != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("terjadi panic saat menambahkan category: %v", r)
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal")
		}
	}()

	rollbackWithError := func(err error) {
		tx.Rollback()
		log.Printf("gagal menambahkan data category: %v", err)
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan data: "+err.Error())
	}

	code := generateCategoryCode()

	var existingCategory category.Category
	if err := tx.Where("name = ?", request.Name).First(&existingCategory).Error; err == nil {
		rollbackWithError(fmt.Errorf("kategori dengan nama tersebut sudah ada"))
		return
	} else if err != gorm.ErrRecordNotFound {
		rollbackWithError(err)
		return
	}

	newCategory := category.Category{
		Code:        code,
		Name:        request.Name,
		Description: request.Description,
		StatusTxt:   "active",
	}

	if err := tx.Create(&newCategory).Error; err != nil {
		rollbackWithError(err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		rollbackWithError(err)
		return
	}

	helpers.SuccessResponse(c, http.StatusCreated, "Data berhasil ditambahkan", newCategory)
}
