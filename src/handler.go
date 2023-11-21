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

type IndexData struct {
	Title       string
	Description string
	RoomName    string
	Modifier    string
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

	d := IndexData{
		Title:       "Challenge you frinds🤼, match more cards🎴 and victory🏅",
		Description: "Name you room, select how many player and have fun!",
		RoomName:    GeneratePetName(3, "-"),
	}

	tmpl.Execute(w, d)
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

	redirect := fmt.Sprintf("/room/%s", path.String())

	http.Redirect(w, r, redirect, http.StatusFound)
}

func (ws *wsHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/room/")

	if _, ok := ws.broker.Rooms[name]; !ok {
		ws.RoomNotFound(w, r)
		return
	}

	room := ws.broker.Rooms[name]
	if len(room.Clients) >= room.ConnLen.Int() {
		ws.RoomAreadyFull(w, r)
		return
	}

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
	cl := NewClient(conn, ch, ws.broker.Rooms[name])

	ws.broker.Register <- cl

	go cl.Write()
	cl.Read(ws.broker)
}

func (ws *wsHandler) RoomNotFound(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/room/")

	tmpl, err := template.ParseFiles("./view/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	d := IndexData{
		Title:       "Room not found! 💡Create this room and play.",
		Description: "select how many player and have fun!",
		RoomName:    name,
		Modifier:    "notfound",
	}

	tmpl.Execute(w, d)
}

func (ws *wsHandler) RoomAreadyFull(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./view/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	d := IndexData{
		Title:       "All player aready playing🎮! Create a new 🌟room and play▶️",
		Description: "Name you room, select how many player and have fun!",
		RoomName:    GeneratePetName(3, "-"),
		Modifier:    "areadyfull",
	}

	tmpl.Execute(w, d)
}
