package src

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type wsHandler struct {
	broker *Broker
}

func NewWsHandler(b *Broker) *wsHandler {
	return &wsHandler{
		broker: b,
	}
}

func (ws *wsHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("./view/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, "")
}

func (ws *wsHandler) RegisterRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	form := r.Form

	path, err := NewPath(form["room"][0])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nplayers, err := NewPlayers(form["players"][0])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := ws.broker.Rooms[path.String()]; ok {
		http.Error(w, "Room aready exists", http.StatusConflict)
		return
	}

	room := NewRoom(path, Players(nplayers))
	ws.broker.Rooms[path.String()] = room
	RoomIDToPath[room.ID] = path.String()

	redirect := fmt.Sprintf("/room/%s", path.String())

	http.Redirect(w, r, redirect, http.StatusFound)
}

func (ws *wsHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/room/")

	if _, ok := ws.broker.Rooms[name]; !ok {
		err := fmt.Sprintf("Room with name %s does not exist", name)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	room := ws.broker.Rooms[name]
	if len(room.Clients) >= room.ConnLen.Int() {
		err := fmt.Sprintf("Room %s is full", name)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// template for room
	tmpl, err := template.ParseFiles("./view/memory.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, ws.broker.Rooms[name])
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *wsHandler) Socket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/ws/room/")

	ch := make(chan *Message, 1)
	cl := NewClient(conn, ch, ws.broker.Rooms[name].ID)

	ws.broker.Register <- cl

	go cl.Write()
	cl.Read(ws.broker)
}
