package main

import (
	"log"

	"mangahub/internal/tcp"
	"mangahub/pkg/database"
)

func main() {
	db, err := database.InitSQLite("data/mangahub.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("[TCP] starting on :9090")

	tcp.StartTCPServer(":9090", db)
}
