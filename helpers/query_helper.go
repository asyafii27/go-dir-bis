package helpers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PreloadIfParam(c *gin.Context, db *gorm.DB, param string, relation string) *gorm.DB {
	if c.Query(param) != "" {
		return db.Preload(relation)
	}
	return db
}
