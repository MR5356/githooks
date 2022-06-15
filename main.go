package main

import (
	"flag"
	"fmt"
	"githooks/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	port  = flag.Int("port", 8900, "port")
	host  = flag.String("host", "0.0.0.0", "host")
	debug = flag.Bool("debug", false, "debug")
)

func RunWeb(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("%v", &json)
	c.JSON(http.StatusOK, 123)
}

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
	r.POST("/", RunWeb)
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
