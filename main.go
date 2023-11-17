package main

import (
	"log"
	"net/http"

	"github.com/tsntt/memory/src"
)

var emojis []string

func main() {

	src.FromJson("emojis.json", &emojis)

	broker := src.NewBroker()
	handlers := src.NewWsHandler(broker)

	http.HandleFunc("/", handlers.Register)
	http.HandleFunc("/room", handlers.RegisterRoom)
	http.HandleFunc("/room/", handlers.JoinRoom)
	http.HandleFunc("/ws/room/", handlers.Socket)

	log.Println("server running at 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
