package controllers

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"slackbot-test/services"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func ProcessEvents(gctx *gin.Context) {
	body, err := ioutil.ReadAll(gctx.Request.Body)
	if err != nil {
		gctx.Error(err)
		return
	}

	sv, err := slack.NewSecretsVerifier(gctx.Request.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		gctx.Error(err)
		return
	}

	if _, err := sv.Write(body); err != nil {
		gctx.Error(err)
		return
	}
	if err := sv.Ensure(); err != nil {
		gctx.Error(err)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		gctx.Error(err)
		return
	}

	services.HandleEvent(gctx, eventsAPIEvent, body)
}

