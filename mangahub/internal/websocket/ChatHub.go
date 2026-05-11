package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type ChatRoom struct {
	ID           string
	MangaID      string
	Participants map[string]*websocket.Conn
	Messages     []ChatMessage
}

type ClientConnection struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
	MangaID  string
}

type RoomMessage struct {
	MangaID string
	Message ChatMessage
}

type ChatHub struct {
	Rooms      map[string]*ChatRoom
	Register   chan ClientConnection
	Unregister chan ClientConnection
	Broadcast  chan RoomMessage
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		Rooms:      make(map[string]*ChatRoom),
		Register:   make(chan ClientConnection),
		Unregister: make(chan ClientConnection),
		Broadcast:  make(chan RoomMessage),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:
			room, exists := h.Rooms[client.MangaID]
			if !exists {
				room = &ChatRoom{
					ID:           client.MangaID,
					MangaID:      client.MangaID,
					Participants: make(map[string]*websocket.Conn),
					Messages:     []ChatMessage{},
				}
				h.Rooms[client.MangaID] = room
				fmt.Println("[WS] Room created:", client.MangaID)
			}

			room.Participants[client.UserID] = client.Conn
			fmt.Printf("[WS] User connected: %s joined room: %s\n", client.Username, client.MangaID)

		case client := <-h.Unregister:
			room, exists := h.Rooms[client.MangaID]
			if exists {
				delete(room.Participants, client.UserID)
				client.Conn.Close()

				fmt.Printf("[WS] User disconnected: %s left room: %s\n", client.Username, client.MangaID)

				if len(room.Participants) == 0 {
					delete(h.Rooms, client.MangaID)
					fmt.Println("[WS] Room removed:", client.MangaID)
				}
			}

		case roomMsg := <-h.Broadcast:
			room, exists := h.Rooms[roomMsg.MangaID]
			if !exists {
				continue
			}

			room.Messages = append(room.Messages, roomMsg.Message)

			for userID, conn := range room.Participants {
				if err := conn.WriteJSON(roomMsg.Message); err != nil {
					fmt.Println("[WS] Write error:", err)
					conn.Close()
					delete(room.Participants, userID)
				}
			}
		}
	}
}
func (h *ChatHub) ActiveUsers() int {
	total := 0

	for _, room := range h.Rooms {
		total += len(room.Participants)
	}

	return total
}
