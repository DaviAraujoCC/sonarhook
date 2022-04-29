package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	// GoogleChatWebhookURL is the URL of the Google Chat webhook
	GoogleChatWebhookURL = os.Getenv("GOOGLE_CHAT_WEBHOOK_URL")
	// ServerPort is the port of the server
	ServerPort = os.Getenv("SERVER_PORT")
	// Status is the filter for results
	Status = os.Getenv("STATUS")
)

func init() {
	switch {
	case GoogleChatWebhookURL == "":
		log.Println("GOOGLE_CHAT_WEBHOOK_URL is not set")
		os.Exit(1)
	case ServerPort == "":
		log.Println("SERVER_PORT is not set, using default 30000")
		ServerPort = "30000"
	}
}
