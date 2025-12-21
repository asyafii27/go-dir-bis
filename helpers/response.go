package helpers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func IsDebug() bool {
	return os.Getenv("APP_DEBUG") == "true"
}

func ErrorResponse(c *gin.Context, code int, customMessage string, err error) {
	message := customMessage

	if IsDebug() && err != nil {
		message = fmt.Sprintf("%s: %s", customMessage, err.Error())
	}

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
