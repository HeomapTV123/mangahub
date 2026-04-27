package tcp

import (
	"database/sql"
	"log"
	"net"
)

func StartTCPServer(address string, db *sql.DB) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("[TCP] failed to start server:", err)
	}

	log.Println("[TCP] Progress Sync Server running on", address)

	hub := NewHub()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[TCP] accept error:", err)
			continue
		}

		go HandleClient(conn, hub, db)
	}
}
