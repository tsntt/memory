package src

import (
	"fmt"
	"log"
)

type Broker struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewBroker() (broker *Broker) {
	broker = &Broker{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 1),
	}

	go broker.listen()

	return
}

func (b *Broker) listen() {
	for {
		select {
		case cl := <-b.Register:
			if _, ok := b.Rooms[cl.RoomPath]; ok {
				r := b.Rooms[cl.RoomPath]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl

					// broadcast message to clients of new user?
					for _, c := range r.Clients {
						if c.ID != cl.ID {
							p := Player{
								ID:           c.ID.String(),
								Username:     c.Username,
								PlayerNumber: c.PlayerNumber,
							}

							strHtml, err := NewUserNameHTMl(p)
							if err != nil {
								log.Print(err)
							}

							cl.Message <- &Message{
								RoomPath: cl.RoomPath,
								UserID:   cl.ID,
								Content:  strHtml,
								Event:    "addUser",
							}
						}
					}
				}
			}
		case cl := <-b.Unregister:
			if _, ok := b.Rooms[cl.RoomPath]; ok {
				r := b.Rooms[cl.RoomPath]

				if _, ok := r.Clients[cl.ID]; ok {
					// if has other users
					if len(r.Clients) > 1 {
						// broadcast User out
						b.Broadcast <- &Message{
							RoomPath: cl.RoomPath,
							UserID:   cl.ID,
							Content:  fmt.Sprintf("%s left the room", cl.Username),
							Event:    "leave",
						}
					}

					delete(r.Clients, cl.ID)
					close(cl.Message)
				}
			}
		case m := <-b.Broadcast:
			boadcast := true

			if _, ok := b.Rooms[m.RoomPath]; ok {
				r := b.Rooms[m.RoomPath]

				// exec things based on m.Action
				switch m.Event {
				case "updateUserName":
					name := fmt.Sprint(m.Content)

					r.Clients[m.UserID].Username = name

					u := Player{
						ID:           m.UserID.String(),
						PlayerNumber: r.Clients[m.UserID].PlayerNumber,
						Username:     name,
					}

					htmlStr, err := NewUserNameHTMl(u)
					if err != nil {
						log.Print(err)
					}

					m.Event = "addUser"
					m.Content = htmlStr
				case "userAdded":
					boadcast = false
					// check this part
					if len(r.Clients) == r.ConnLen.Int() {
						boadcast = true

						m.Event = "start"
						m.Content = map[string]*Client{
							"turn": r.Clients[r.ClientTurn()],
						}
					}

				}

				// condition to send or not
				if boadcast {
					for _, cl := range r.Clients {
						cl.Message <- m
					}
				}
			}
		}
	}
}
