package hooks

import (
	"encoding/json"
	"githooks/config"
	"githooks/runner"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/gitlab"
	"log"
	"net/http"
)

func HandleGitlab(c *gin.Context) {
	hook, _ := gitlab.New(gitlab.Options.Secret(config.Secret))

	payload, err := hook.Parse(c.Request, gitlab.PushEvents)
	if err != nil {
		log.Println(err)
	}

	payloadJson, _ := json.Marshal(payload)
	log.Printf("new gitlab hook: %+v", string(payloadJson))

	switch payload.(type) {
	case gitlab.PushEventPayload:
		pl := payload.(gitlab.PushEventPayload)

		builder := runner.NewDefaultBuild()
		builder.Name = pl.Project.Name
		builder.From = "Gitlab"
		builder.SshUrl = pl.Project.GitSSHURL
		builder.HttpUrl = pl.Project.GitHTTPURL
		builder.CommitId = pl.After[0:6]
		builder.UserName = pl.UserName
		go builder.Run()

		//go utils.RunScript("docker.sh", []string{pl.Project.Name, pl.Project.GitSSHURL, pl.After[0:6]})
	}

	c.JSON(http.StatusOK, payload)
}
