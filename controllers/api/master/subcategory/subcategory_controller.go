package subcategory

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/master/subcategory"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSubCategories(c *gin.Context) {
	var subcategories []subcategory.SubCategory

	db := config.Database

	db = db.Preload("Category")

	db = ApplySubCategoryFilters(c, db)

	page := c.Query("page")

	if page != "" {
		meta, err := helpers.LaravelPaginate(c, db, &subcategories)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, meta)
		return
	}

	if err := db.Find(&subcategories).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subcategories)

}

func ApplySubCategoryFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
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

func GetSubCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var subcategory subcategory.SubCategory

	db := config.Database

	db = db.Preload("Category")

	if err := db.First(&subcategory, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "Subcategory tidak ditemukan")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", subcategory)
}
