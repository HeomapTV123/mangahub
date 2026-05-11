package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	token := flag.String("token", "", "JWT token from login")
	mangaID := flag.String("manga_id", "", "manga room id")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required. Example: go run ./cmd/ws-client -token=YOUR_TOKEN -manga_id=one-piece")
	}

	if *mangaID == "" {
		log.Fatal("manga_id is required. Example: go run ./cmd/ws-client -token=YOUR_TOKEN -manga_id=one-piece")
	}

	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws",
	}

	query := u.Query()
	query.Set("token", *token)
	query.Set("manga_id", *mangaID)
	u.RawQuery = query.Encode()

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if resp != nil {
			log.Fatalf("failed to connect websocket: %v, HTTP status: %s", err, resp.Status)
		}
		log.Fatal("failed to connect websocket: ", err)
	}
	defer conn.Close()

	fmt.Println("Connected to MangaHub WebSocket room:", *mangaID)
	fmt.Println("Type a message and press Enter.")
	fmt.Println("Type /quit to exit.")

	go func() {
		for {
			var msg ChatMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("Disconnected from server.")
				os.Exit(0)
			}

			fmt.Printf("\n[%s | %s]: %s\n> ", msg.MangaID, msg.Username, msg.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()

		if text == "/quit" {
			fmt.Println("Closing chat...")
			return
		}

		msg := ChatMessage{
			Message:   text,
			Timestamp: time.Now().Unix(),
		}

		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("send error:", err)
			return
		}

		fmt.Print("> ")
	}
}
