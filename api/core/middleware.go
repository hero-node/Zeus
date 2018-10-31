package core

import (
	"github.com/gin-gonic/gin"
)

func IpfsAddFilter() gin.HandlerFunc {
	return func(c *gin.Context) {

		if true {
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"result": "error",
				"reason": "Not enought HER",
			})
			c.Abort()
		}
	}
}
