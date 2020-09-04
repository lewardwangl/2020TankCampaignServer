package utils

import (
	"io"
	"log"
	"os"
	"server/config"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Request *log.Logger
)

func init() {
	errFile, err := OpenFile(config.LogStoreDir + "errors.log")
	info, err := OpenFile(config.LogStoreDir + "info.log")
	warning, err := OpenFile(config.LogStoreDir + "warning.log")
	request, err := OpenFile(config.LogStoreDir + "req.log")
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	Info = log.New(io.MultiWriter(os.Stdout, info), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stdout, warning), "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	Request = log.New(io.MultiWriter(os.Stderr, request), "Request:", log.Ldate|log.Ltime|log.Lshortfile)
}
