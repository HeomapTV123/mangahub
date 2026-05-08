package main

import (
	"mangahub/internal/udp"
)

func main() {
	server := &udp.NotificationServer{
		Port: ":9091",
	}

	udp.StartUDPServer(server)
}
