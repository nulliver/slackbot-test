package main

import (
	"net/http"

	"slackbot-test/controllers"
	"slackbot-test/logger"
	"slackbot-test/storage"
)

func main() {


	http.HandleFunc("/slack/events", controllers.ProcessEvent)

	storage.Setup()

	logger.Info("Nerdcoin bot started")
	http.ListenAndServe(":3000", nil)
}