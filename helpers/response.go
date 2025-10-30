package helpers

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  code,
		"message": message,
	})
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}, meta ...interface{}) {
	response := gin.H{
		"status":  code,
		"message": message,
		"data":    data,
	}

	if len(meta) > 0 && meta[0] != nil {
		response["meta"] = meta[0]
	}

	c.JSON(code, response)
}
