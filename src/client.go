package src

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn `json:"-"`
	Message  chan *Message   `json:"-"`
	ID       uuid.UUID       `json:"id"`
	RoomID   uuid.UUID       `json:"-"`
	Username string          `json:"username"`
}

type Message struct {
	RoomID   uuid.UUID `json:"roomId"`
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Action   string    `json:"action"`
}

func NewClient(conn *websocket.Conn, ch chan *Message, roomID uuid.UUID) *Client {
	return &Client{
		ID:      uuid.New(),
		Conn:    conn,
		Message: ch,
		RoomID:  roomID,
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

		msg.RoomID = c.RoomID
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
