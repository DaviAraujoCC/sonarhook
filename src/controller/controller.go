package controller

import (
	"encoding/json"
	"net/http"
	"sonarhook/src/message"

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

	// Message constructor
	mc := message.NewMessage(msg)

	// Send the message
	err = mc.SendMessage()
	if err != nil {
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
