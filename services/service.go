package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"slackbot-test/storage"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))

func HandleEvent(gctx *gin.Context, eventsAPIEvent slackevents.EventsAPIEvent, body []byte) {
	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		handleUrlVerification(gctx, body)
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			processSlackMessage(ev)
		}

	}
}

func processSlackMessage(ev *slackevents.MessageEvent) {

	if strings.Contains(ev.Text, "++") {

		var usersWithPlusPlus []string
		plusPlusRegex := regexp.MustCompile(`@[^\s@+]*?\+{2}`)
		matches := plusPlusRegex.FindAllString(ev.Text, -1)
		for _, userId := range matches {
			log.Print("userId (before trim): " + userId)
			userId = strings.TrimPrefix(userId, "@")
			userId = strings.TrimSuffix(userId, ">++")
			log.Print("userId (after trim): " + userId)
			userInfo, err := api.GetUserInfo(userId)
			if err != nil {
				log.Print(err)
				continue
			}
			usersWithPlusPlus = append(usersWithPlusPlus, userInfo.Name)
		}
		api.PostMessage(ev.Channel, slack.MsgOptionText("Coins :nerdcoin: for: " + strings.Join(usersWithPlusPlus, ", "), false), slack.MsgOptionTS(ev.Message.TimeStamp))
		log.Printf("ev.User: %s", ev.User)
		userInfo, err := api.GetUserInfo(ev.User)
		if err != nil {
			log.Print(err)
		}

		storage.SaveTransaction(userInfo.Name, ev.Text, usersWithPlusPlus)
	}
}

func handleUrlVerification(gctx *gin.Context, body []byte) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		gctx.Error(err)
		return
	}
	gctx.JSON(http.StatusOK, gin.H{"challenge": r.Challenge})
	/*w.Header().Set("Content-Type", "text")
	w.Write([]byte(r.Challenge))*/
}
