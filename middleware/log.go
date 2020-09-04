package middleware

import (
	"github.com/gin-gonic/gin"
	"server/config"
	"server/utils"
	"strconv"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		utils.Request.Println("["+strconv.Itoa(int(time.Since(now).Milliseconds()))+"ms]",
			"["+c.Request.Method+"]", c.ClientIP(), c.Request.URL, c.Writer.Status(),
			"__"+c.GetString(config.RequestLogParamKey)+"__",
		)
	}
}
