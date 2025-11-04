package partnerowner

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/partnerowner"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPartnerOwners(c *gin.Context) {
	var partnerowners []partnerowner.PartnerOwner
	db := config.Database

	db = db.Preload("Province").Preload("City").Preload("District").Preload("Village")

	db = ApplyPartnerOwnerFilters(c, db)

	page := c.Query("page")

	if page != "" {
		meta, err := helpers.LaravelPaginate(c, db, &partnerowners)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, meta)
		return
	}

	if err := db.Find(&partnerowners).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, partnerowners)
}

func ApplyPartnerOwnerFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if globalSearch := c.Query("global_search"); globalSearch != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+globalSearch+"%", "%"+globalSearch+"%")
	}

	if address := c.Query("address"); address != "" {
		db = db.Where("address = ?", address)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":         true,
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

func GetPartnerOwnerByID(c *gin.Context) {
	id := c.Param("id")
	var partnerowner partnerowner.PartnerOwner

	db := config.Database

	db = db.Preload("Province").Preload("City").Preload("District").Preload("Village")

	if err := db.First(&partnerowner, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "Data pemilik tidak ditemukan")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", partnerowner)
}
