package main

import (
	"flag"
	"fmt"
	"githooks/hooks"
	"githooks/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

var (
	port  = flag.Int("port", 8900, "port")
	host  = flag.String("host", "0.0.0.0", "host")
	debug = flag.Bool("debug", false, "debug")
)

func Init() {
	// 命令行参数解析
	flag.Parse()

	// 命令行清空
	utils.Clear()

	// 日志相关
	logPath := "./logs"
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(logPath, 0644)
		if err != nil {
			fmt.Println("日志目录创建失败")
		}
	}
	logFile, err := os.OpenFile(logPath+"/githooks.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetPrefix("[githooks] ")
	log.SetOutput(mw)
	gin.DefaultWriter = mw
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
}

func CreateRoute() *gin.Engine {
	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.POST("/github", hooks.HandleGithub)
	r.POST("/gitlab", hooks.HandleGitlab)
	// 输出程序端口信息
	log.Printf("Listening and serving HTTP on %s:%d with PID %d", *host, *port, os.Getpid())
	return r
}

func main() {
	// 初始化操作
	Init()

	r := CreateRoute()

	// 保存PID
	utils.RunCommand(fmt.Sprintf("echo %d > run.pid", os.Getpid()))
	err := r.Run(fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		panic(err)
	}
}
