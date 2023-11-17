package src

import (
	"fmt"
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
		Broadcast:  make(chan *Message, 2),
	}

	go broker.listen()

	return
}

func (b *Broker) listen() {
	for {
		select {
		case cl := <-b.Register:
			path := RoomIDToPath[cl.RoomID]

			if _, ok := b.Rooms[path]; ok {
				r := b.Rooms[path]

				if _, ok := r.Clients[cl.ID.String()]; !ok {
					r.Clients[cl.ID.String()] = cl
				}
			}
		case cl := <-b.Unregister:
			path := RoomIDToPath[cl.RoomID]

			if _, ok := b.Rooms[path]; ok {
				r := b.Rooms[path]

				if _, ok := r.Clients[cl.ID.String()]; ok {
					// if has other users
					if len(r.Clients) > 1 {
						// broadcast User out
						b.Broadcast <- &Message{
							RoomID:   cl.RoomID,
							UserID:   cl.ID,
							Username: cl.Username,
							Content:  fmt.Sprintf("%s left the room", cl.Username),
							Action:   "leave",
						}
					}

					delete(r.Clients, cl.ID.String())
					close(cl.Message)
				}
			}
		case m := <-b.Broadcast:
			path := RoomIDToPath[m.RoomID]
			if _, ok := b.Rooms[path]; ok {
				r := b.Rooms[path]

				// exec things based on m.Action
				switch m.Action {
				case "register":

				case "updatename":
					r.Clients[m.UserID.String()].Username = m.Username
					m.Content = ""
					if len(r.Clients) == r.ConnLen.Int() {
						m.Action = "start"
					} else {
						m.Action = "wait"
					}
				}

				for _, cl := range r.Clients {
					cl.Message <- m
				}
			}
		}
	}
}
