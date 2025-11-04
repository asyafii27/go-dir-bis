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

func GetVillages(c *gin.Context) {
	var villages []region.Village

	db := config.Database

	db = ApplyVillageFilters(c, db)
	meta, err := PaginateData(c, db, &villages)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", villages, meta)
}

func ApplyVillageFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if code := c.Query("code"); code != "" {
		db = db.Where("code = ?", code)
	}

	if village_code := c.Query("village_code"); village_code != "" {
		db = db.Where("village_code = ?", village_code)
	}

	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":           true,
		"code":         true,
		"village_code": true,
		"name":         true,
		"created_at":   true,
		"updated_at":   true,
	}

	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return db.Order(sortBy + " " + sortDir)
}

func GetVillageByID(c *gin.Context) {
	id := c.Param("id")
	var village region.Village

	if err := config.Database.First(&village, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "VIllage not found")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", village)
}
