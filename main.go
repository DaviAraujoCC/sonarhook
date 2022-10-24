package main

import (
	"net/http"
	"sonarhook/config"
	"sonarhook/controller"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/webhook", controller.HandleWebhook).Methods("POST")

	server := &http.Server{
		Addr:    ":" + config.ServerPort,
		Handler: r,
	}

	log.Info("Starting server on address " + server.Addr)
	log.Fatal(server.ListenAndServe())

}
