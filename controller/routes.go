package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var err error

func HandleWebhook(w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var msg Message
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Error(err)
		return
	}

	// Validate the message
	err = msg.validateMessage()
	if err != nil {
		log.Error(err)
		return
	}

	// Send the message
	err = msg.sendMessage()
	if err != nil {
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
