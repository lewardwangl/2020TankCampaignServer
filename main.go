package main

import (
	"github.com/gin-gonic/gin"
	"server/handlers/ranks"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.GET("/api/ranks", ranks.GetData)
	r.Any("/api/ws", ranks.WebSocketRanks)
	r.Static("/public", "./public")
	r.Run()
}
