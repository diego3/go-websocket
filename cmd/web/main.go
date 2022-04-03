package main

import (
	"log"
	"net/http"

	"github.com/diego3/go-websocket/internal/handlers"
)

func main() {
	routes()

	log.Println("Starting web server on port 8080")

	go handlers.ListenPayloadChannel()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error trying to server", err)
	}
}
