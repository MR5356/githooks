package hooks

import (
	"encoding/json"
	"fmt"
	"githooks/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type GithubHookBody struct {
	Repository struct {
		Name   string `json:"name"`
		SSHUrl string `json:"ssh_url"`
	} `json:"repository"`
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

	req, _ := json.Marshal(githubHookBody)
	fmt.Printf("接收到新的githook：%s", string(req))

	utils.RunScript("docker.sh", []string{githubHookBody.Repository.Name, githubHookBody.Repository.SSHUrl, githubHookBody.HeadCommit.Id[0:6]})
	c.JSON(http.StatusOK, githubHookBody)
}
