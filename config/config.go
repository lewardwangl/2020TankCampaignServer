package config

import "os"

// some config
var MysqlHost string
var MysqlPort = "3306"
var MysqlUser = "root"
var MysqlPassword string
var MysqlDatabase string
var ListenPort = "8080"
var IsProd = false
var RequestLogParamKey = "_log_param_"

// 验证字符串，写死的 也没登录  就直接比对这个字符串就行了
var AuthSKey = "_tc"
var AuthSKeyValue = "sZXczvbznhjasdzdddasfgddsZXcvddbnhjhtyredwasDFGHJdYTfREaWAsa"

// 日志文件存放地址
var LogStoreDir string

func init() {
	MysqlHost = os.Getenv("MysqlHost")
	MysqlDatabase = os.Getenv("MysqlDatabase")
	MysqlPassword = os.Getenv("MysqlPassword")
	t := os.Getenv("MysqlPort")
	if t != "" {
		MysqlPort = t
	}
	t = os.Getenv("MysqlUser")
	if t != "" {
		MysqlUser = t
	}
	t = os.Getenv("ListenPort")
	if t != "" {
		ListenPort = t
	}
	t = os.Getenv("IsProd")
	if t == "true" {
		IsProd = true
	}
	if IsProd {
		LogStoreDir = "/var/log/1024/"
	} else {
		LogStoreDir = "log/"
	}
}
