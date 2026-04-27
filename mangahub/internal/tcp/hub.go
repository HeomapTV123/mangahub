package tcp

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

type Hub struct {
	clients map[net.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[net.Conn]bool),
	}
}

func (h *Hub) AddClient(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = true
	log.Println("[TCP] client connected:", conn.RemoteAddr())
}

func (h *Hub) RemoveClient(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
	conn.Close()
	log.Println("[TCP] client disconnected:", conn.RemoteAddr())
}

func (h *Hub) Broadcast(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("[TCP] marshal error:", err)
		return
	}

	data = append(data, '\n')

	for conn := range h.clients {
		if _, err := conn.Write(data); err != nil {
			log.Println("[TCP] broadcast error:", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}
