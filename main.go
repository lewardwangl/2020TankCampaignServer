package main

import (
	"github.com/gin-gonic/gin"
	"server/config"
	_ "server/data"
	"server/handlers/leetcode"
	"server/handlers/tank"
	"server/middleware"
	"server/utils"
	_ "server/utils"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.Logger())

	api := r.Group("api")
	{
		// 力扣大赛API
		leetcodeApi := api.Group("leetcode")
		leetcodeApi.GET("ranks", leetcode.Ranks)
	}
	{
		// 坦克大战API
		tankApi := api.Group("tank")
		// 获取车轮战排名
		tankApi.GET("ranks", tank.Ranks)
		// 获取决赛排名
		tankApi.GET("ranks/final", tank.FinalRanks)
		// 查询比赛是否存在 通过battle_id
		tankApi.GET("battle/isexists", middleware.AuthCookie(), tank.HavingBattleId)
		// 插入一场车轮战比赛信息
		tankApi.POST("rank", middleware.AuthCookie(), tank.CreateBattleInfo)
		// 宣布一场比赛无效
		tankApi.PATCH("rank/invalid", middleware.AuthCookie(), tank.InvalidBattleById)
		// 插入一场决赛比赛信息
		tankApi.POST("rank/final", middleware.AuthCookie(), tank.CreateFinalBattleInfo)
		// 实时推送接口
		tankApi.Any("ws", tank.WebSocketRanks)
	}

	r.Static("public", "./public") // test

	err := r.Run(":" + config.ListenPort)
	if err != nil {
		utils.Error.Println(err)
		panic(err)
	}
}
