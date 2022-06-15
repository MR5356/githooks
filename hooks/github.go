package hooks

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type GithubHookBody struct {
	Depository struct {
		Name   string `json:"name"`
		SSHUrl string `json:"ssh_url"`
	} `json:"depository"`
	HeadCommit struct {
		Id string `json:"id"`
	} `json:"head_commit"`
}

func HandleGithub(c *gin.Context) {
	githubHookBody := GithubHookBody{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &githubHookBody)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, githubHookBody)
}
