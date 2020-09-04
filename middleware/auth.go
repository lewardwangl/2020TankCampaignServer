package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/config"
)

func AuthCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query(config.AuthSKey) != config.AuthSKeyValue {
			c.JSON(http.StatusBadRequest, gin.H{"error": "你猜你做错了啥？小样儿"})
			c.Abort()
			return
		}
		c.Next()
	}
}
