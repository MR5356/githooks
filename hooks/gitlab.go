package hooks

import (
	"githooks/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/gitlab"
	"log"
)

func HandleGitlab(c *gin.Context) {
	hook, _ := gitlab.New(gitlab.Options.Secret(config.Secret))

	payload, err := hook.Parse(c.Request, gitlab.PushEvents)
	if err != nil {
		log.Println(err)
	}

	switch payload.(type) {
	case gitlab.PushEventPayload:
		log.Printf("new gitlab hook: %+v", payload)
	}
}
