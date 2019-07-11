package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/novemberde/go-websocket-demo/api"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	hub := api.NewHub()
	go hub.Run()
	// http.HandleFunc("/", serverHome)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
