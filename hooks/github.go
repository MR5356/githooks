package hooks

import (
	"encoding/json"
	"githooks/config"
	"githooks/runner"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/github"
	"log"
	"net/http"
)

func HandleGithub(c *gin.Context) {
	hook, _ := github.New(github.Options.Secret(config.Secret))

	payload, err := hook.Parse(c.Request, github.PushEvent)
	if err != nil {
		log.Println(err)
	}

	payloadJson, _ := json.Marshal(payload)
	log.Printf("new github hook: %+v", string(payloadJson))

	switch payload.(type) {
	case github.PushPayload:
		pl := payload.(github.PushPayload)

		builder := runner.NewDefaultBuild()
		builder.Name = pl.Repository.Name
		builder.From = "Github"
		builder.SshUrl = pl.Repository.SSHURL
		builder.HttpUrl = pl.Repository.CloneURL
		builder.CommitId = pl.After[0:6]
		builder.UserName = pl.Pusher.Name
		go builder.Run()

		//go utils.RunScript("docker.sh", []string{pl.Repository.Name, pl.Repository.SSHURL, pl.After[0:6]})
	}

	c.JSON(http.StatusOK, payload)
}
