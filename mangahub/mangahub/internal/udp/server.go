package udp

import (
	"fmt"
	"net"
)

func StartUDPServer(server *NotificationServer) {
	addr, err := net.ResolveUDPAddr("udp", server.Port)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	server.Conn = conn

	fmt.Println("UDP Server listening on", server.Port)

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		message := string(buffer[:n])

		fmt.Printf("Received from %s: %s\n", clientAddr, message)

		if message == "register" {
			registerClient(server, clientAddr)

			conn.WriteToUDP([]byte("registered"), clientAddr)
		}
	}
}
