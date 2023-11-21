package src

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn         *websocket.Conn `json:"-"`
	Message      chan *Message   `json:"-"`
	ID           uuid.UUID       `json:"id"`
	RoomPath     string          `json:"-"`
	Username     string          `json:"username"`
	PlayerNumber int             `json:"playerNumber"`
}

type Message struct {
	RoomPath string    `json:"roomId"`
	UserID   uuid.UUID `json:"userId"`
	Content  any       `json:"content"`
	Event    string    `json:"event"`
}

func NewClient(conn *websocket.Conn, ch chan *Message, room *Room) *Client {
	return &Client{
		ID:           uuid.New(),
		Conn:         conn,
		Message:      ch,
		RoomPath:     room.Path.String(),
		PlayerNumber: len(room.Clients) + 1,
	}
}

func (c *Client) Read(broker *Broker) {
	defer func() {
		broker.Unregister <- c
		c.Conn.Close()
	}()

	for {

		msg := Message{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg.RoomPath = c.RoomPath
		msg.UserID = c.ID

		broker.Broadcast <- &msg
	}
}

func (c *Client) Write() {
	defer c.Conn.Close()

	for {
		m, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(m)
	}

}
