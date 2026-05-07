package udp

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func registerClient(server *NotificationServer, addr *net.UDPAddr) {
	for _, c := range server.Clients {
		if c.String() == addr.String() {
			return // already registered
		}
	}

	server.Clients = append(server.Clients, *addr)
	fmt.Println("UDP Server Client registered:", addr)
}
func broadcast(server *NotificationServer, conn *net.UDPConn, notif Notification) {
	data, err := json.Marshal(notif)
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	for _, client := range server.Clients {
		_, err := conn.WriteToUDP(data, &client)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}
}
func broadcastLoop(server *NotificationServer, conn *net.UDPConn) {
	for {
		time.Sleep(10 * time.Second)

		notif := Notification{
			Type:      "new_chapter",
			MangaID:   "naruto",
			Message:   "Naruto Chapter 700 released!",
			Timestamp: time.Now().Unix(),
		}

		broadcast(server, conn, notif)

		fmt.Println("[UDP] Broadcast sent to", len(server.Clients), "clients")
	}
}
