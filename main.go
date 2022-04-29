package main

import (
	"net/http"
	"sonarhook/controller"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/webhook", controller.HandleWebhook).Methods("POST")

	server := &http.Server{
		Addr:    ":30000",
		Handler: r,
	}

	log.Info("Starting server on port " + server.Addr)
	log.Fatal(server.ListenAndServe())

}
