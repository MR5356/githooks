package main

import (
	"flag"
	"fmt"
	"githooks/hooks"
	"githooks/utils"
	"github.com/gin-gonic/gin"
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

	// 输出程序端口信息
	fmt.Printf("Listening and serving HTTP on %s:%d\n", *host, *port)
}

func CreateRoute() *gin.Engine {
	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.POST("/", hooks.HandleGithub)
	return r
}

func main() {
	Init()

	r := CreateRoute()
	err := r.Run(fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		panic(err)
	}
}
