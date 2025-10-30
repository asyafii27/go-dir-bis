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

func GetProvinces(c *gin.Context) {
	var provinces []region.Province
	db := config.Database

	db = ApplyProvinceFilters(c, db)

	meta, err := PaginateData(c, db, &provinces)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "success", provinces, meta)
}

func ApplyProvinceFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if code := c.Query("code"); code != "" {
		db = db.Where("code = ?", code)
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

func PaginateData(c *gin.Context, db *gorm.DB, result interface{}) (helpers.PaginationMeta, error) {
	meta, err := helpers.Paginate(c, db, result)
	if err != nil {
		return meta, err
	}
	return meta, nil
}
