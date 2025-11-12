package secondsubcategory

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/master/secondsubcategory"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSecondSubCategories(c *gin.Context) {
	var secondsubcategories []secondsubcategory.SecondSubCategory

	db := config.Database

	db = db.Preload("Category").Preload("SubCategory.Category")

	db = ApplySecondSubCategoryFilters(c, db)

	page := c.Query("page")

	if page != "" {
		meta, err := helpers.LaravelPaginate(c, db, &secondsubcategories)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, meta)
		return
	}

	if err := db.Find(&secondsubcategories).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, secondsubcategories)

}

func ApplySecondSubCategoryFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	if subCategoryID := c.Query("sub_category_id"); subCategoryID != "" {
		db = db.Where("sub_category_id = ?", subCategoryID)
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

func GetSecondSubCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var secondsubcategory secondsubcategory.SecondSubCategory

	db := config.Database

	db = db.Preload("Category").Preload("SubCategory")

	if err := db.First(&secondsubcategory, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "Second sub Category tidak ditemukan")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", secondsubcategory)
}
