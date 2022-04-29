package controller

import (
	"encoding/json"
	"net/http"
	"sonarhook/message"

	log "github.com/sirupsen/logrus"
)

var err error

func HandleWebhook(w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var msg message.Message
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Error(err)
		return
	}

	// Validate the message
	err = msg.ValidateMessage()
	if err != nil {
		log.Error(err)
		return
	}

	// Send the message
	err = msg.SendMessage()
	if err != nil {
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
