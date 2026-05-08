package udp

import (
	"encoding/json"
	"fmt"
	"net"
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
func Broadcast(server *NotificationServer, notif Notification) {
	data, err := json.Marshal(notif)
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	for _, client := range server.Clients {
		_, err := server.Conn.WriteToUDP(data, &client)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	fmt.Println("[UDP] Notification broadcasted")
}
