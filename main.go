package main

import (
	"flag"
	"fmt"
	"githooks/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	p = flag.Int("p", 8900, "port")
	h = flag.String("h", "0.0.0.0", "host")
)

func RunWeb(c *gin.Context) {
	c.JSON(http.StatusOK, "123")
}

func main() {
	flag.Parse()
	utils.Clear()
	fmt.Println(utils.GetExtFiles(utils.GetAbsPath(), ".exe"))
	q := utils.RunCommand("cmd /c ping 43.138.31.224")
	fmt.Println(q)
	fmt.Printf("Listening and serving HTTP on %s:%d\n", *h, *p)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", RunWeb)
	err := r.Run(fmt.Sprintf("%s:%d", *h, *p))
	if err != nil {
		panic(err)
	}
}
