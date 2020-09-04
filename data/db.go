package db

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"server/config"
	"server/utils"
	"time"
)

var DB *gorm.DB

func init() {
	sqlLogger, err := utils.OpenFile(config.LogStoreDir + "sql.log")
	if err != nil {
		log.Println(err.Error())
		panic("sql日志文件打开失败")
	}
	var writer io.Writer
	if config.IsProd {
		writer = sqlLogger
	} else {
		writer = io.MultiWriter(os.Stdout, sqlLogger)
	}
	logger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // 禁用彩色打印
		},
	)
	DB, err = gorm.Open(mysql.Open(config.MysqlUser+":"+
		config.MysqlPassword+"@tcp("+config.MysqlHost+":"+config.MysqlPort+")/"+
		config.MysqlDatabase), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
}
