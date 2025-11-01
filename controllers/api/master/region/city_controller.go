package region

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/master/region"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCities(c *gin.Context) {
	var cities []region.City
	db := config.Database

	db = AppyCityFilters(c, db)

	meta, err := PaginateData(c, db, &cities)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", cities, meta)
}

func AppyCityFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if code := c.Query("code"); code != "" {
		db = db.Where("code = ?", code)
	}

	if province_code := c.Query("province_code"); province_code != "" {
		db = db.Where("province_code = ?", province_code)
	}

	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":         true,
		"code":       true,
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

func GetCityByID(c *gin.Context) {
	id := c.Param("id")
	var city region.City

	if err := config.Database.First(&city, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "Data tidak ditemukan")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "SUccess", city)
}
