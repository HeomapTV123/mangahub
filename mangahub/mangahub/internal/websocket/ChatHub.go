package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ChatHub struct {
	Clients    map[*websocket.Conn]string
	Broadcast  chan ChatMessage
	Register   chan ClientConnection
	Unregister chan *websocket.Conn
}

type ClientConnection struct {
	Conn     *websocket.Conn
	Username string
}

type ChatMessage struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Conn] = client.Username
			fmt.Println("[WS] User connected:", client.Username)

		case conn := <-h.Unregister:
			delete(h.Clients, conn)
			conn.Close()
			fmt.Println("[WS] User disconnected")

		case msg := <-h.Broadcast:
			for conn := range h.Clients {
				conn.WriteJSON(msg)
			}
		}
	}
}
