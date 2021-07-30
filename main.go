package main

import (
	"log"
	"net/http"
	"os"

	"slackbot-test/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("NRDPORT")
	if port == "" {
		log.Fatal("NRDPORT must be set")
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "Nerdcoin Bot", "message": "Nerdcoin Bot is up and running"})

	})
	// router.POST("/slack/events", controllers.ProcessEvents)

	storage.Setup()

	router.Run(":" + port)
}