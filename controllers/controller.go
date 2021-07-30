package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"slackbot-test/services"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func ProcessEvents(gctx *gin.Context) {
	body, err := ioutil.ReadAll(gctx.Request.Body)
	if err != nil {
		log.Print(err)
		gctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	sv, err := slack.NewSecretsVerifier(gctx.Request.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Print(err)
		gctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := sv.Write(body); err != nil {
		log.Print(err)
		gctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := sv.Ensure(); err != nil {
		log.Print(err)
		gctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Print(err)
		gctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	services.HandleEvent(gctx, eventsAPIEvent, body)
}

