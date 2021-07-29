package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))
var signingSecret   = os.Getenv("SLACK_SIGNING_SECRET")

func ProcessEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handleEvent(w, eventsAPIEvent, body)
}

func handleEvent(w http.ResponseWriter, eventsAPIEvent slackevents.EventsAPIEvent, body []byte) {
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
	}
}

func handleUrlVerification(w http.ResponseWriter, body []byte) bool {
	fmt.Println("- URL Verification")
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