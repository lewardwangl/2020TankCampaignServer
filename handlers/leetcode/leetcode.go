package leetcode

import (
	"github.com/gin-gonic/gin"
	"net/http"
	db "server/data"
)

// Ranks 获取力扣大赛的排名情况 这个是手动操作的数据库数据，数据库数据已经是完好的数据，直接获取
func Ranks(c *gin.Context) {
	var result = []db.LeetCode{}
	db.DB.Debug().Model(&db.LeetCode{}).Order("rank").Find(&result)
	c.JSON(http.StatusOK, result)
}
