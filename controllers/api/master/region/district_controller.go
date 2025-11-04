package region

import (
	"net/http"
	"strings"

	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/master/region"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDistricts(c *gin.Context) {
	var districts []region.District

	db := config.Database

	db = ApplyDistrictFilters(c, db)

	meta, err := PaginateData(c, db, &districts)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", districts, meta)
}

func ApplyDistrictFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if code := c.Query("code"); code != "" {
		db = db.Where("code = ?", code)
	}

	if city_code := c.Query("city_code"); city_code != "" {
		db = db.Where("city_code = ?", city_code)
	}

	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":         true,
		"code":       true,
		"city_code":  true,
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return db.Order(sortBy + " " + sortDir)
}

func GetDistrictByID(c *gin.Context) {
	id := c.Param("id")
	var district region.District

	if err := config.Database.First(&district, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "District not found")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", district)
}
