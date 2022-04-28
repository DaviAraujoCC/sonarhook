package main

import (
	"encoding/json"

	"net/http"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var err error

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/webhook", processMessage).Methods("POST")

	server := &http.Server{
		Addr:    ":30000",
		Handler: r,
	}

	log.Info("Starting server on port " + server.Addr)
	log.Fatal(server.ListenAndServe())

}

func processMessage(w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var msg Message
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Error(err)
		return
	}

	// Get the message
	err = msg.validateMessage()
	if err != nil {
		log.Error(err)
	}
	err = msg.sendMessage()
	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
