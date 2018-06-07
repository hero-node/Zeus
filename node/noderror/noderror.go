package noderror

import (
	"github.com/gin-gonic/gin"
)

// Error : error response
func Error(err error, c *gin.Context) {
	c.JSON(500, gin.H{
		"status": "error",
		"reason": err.Error(),
	})
}
