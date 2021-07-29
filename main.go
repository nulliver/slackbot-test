package main

import (
	"fmt"
	"net/http"

	"slackbot-test/controllers"
)

func main() {
	http.HandleFunc("/slack/events", controllers.ProcessEvent)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}