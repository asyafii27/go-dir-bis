package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Paginate(c *gin.Context, db *gorm.DB, out interface{}) (PaginationMeta, error) {
	var meta PaginationMeta

	pageStr := c.Query("page")
	if pageStr == "" {
		if err := db.Find(out).Error; err != nil {
			return meta, err
		}
		return meta, nil
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	if err := db.Model(out).Count(&total).Error; err != nil {
		return meta, err
	}

	if err := db.Limit(limit).Offset(offset).Find(out).Error; err != nil {
		return meta, err
	}

	meta = PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	return meta, nil
}
