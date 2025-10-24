package helpers

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  code,
		"message": message,
		"data":    data,
	})
}
