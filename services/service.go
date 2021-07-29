package services

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"slackbot-test/storage"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.mongodb.org/mongo-driver/bson"
)

var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))

func HandleEvent(w http.ResponseWriter, eventsAPIEvent slackevents.EventsAPIEvent, body []byte) {
	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		handleUrlVerification(w, body)
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
		for _, user := range matches {
			user = strings.TrimPrefix(user, "@")
			user = strings.TrimSuffix(user, "++")
			usersWithPlusPlus = append(usersWithPlusPlus, user)
		}
		users := strings.Join(usersWithPlusPlus, " ")
		api.PostMessage(ev.Channel, slack.MsgOptionText("Users with plus-plus: " + users, false))

		var docs []interface{}

		for _, u := range usersWithPlusPlus {
			doc := bson.D{
				{Key: "fromUser", Value: ev.User},
				{Key: "toUser", Value: u},
				{Key: "message", Value: ev.Message},
				{Key: "type", Value: "karma"},
				{Key: "timestamp", Value: time.Now()},
			}
			docs = append(docs, doc)
		}
		storage.Db.InsertMany(storage.Ctx, docs)
	}
}

func handleUrlVerification(w http.ResponseWriter, body []byte) bool {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	w.Header().Set("Content-Type", "text")
	w.Write([]byte(r.Challenge))
	return false
}
