package main

import (
	"net/http"

	"github.com/diego3/go-websocket/internal/handlers"
)

func routes() /*http.Handler */ {
	//mux := pat.New()

	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", fs)
	//http.Handle("/", http.StripPrefix("/html/", fs))
	http.Handle("/ws", http.HandlerFunc(handlers.WsEndpoint))

	//return mux
}
