package tcp

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net"

	"mangahub/internal/auth"
)

func HandleClient(conn net.Conn, hub *Hub, db *sql.DB) {
	hub.AddClient(conn)
	defer hub.RemoveClient(conn)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		var msg Message

		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Println("[TCP] invalid JSON:", err)
			send(conn, Message{Type: "error", Text: "invalid json"})
			continue
		}

		switch msg.Type {
		case "ping":
			send(conn, Message{Type: "pong", Text: "server alive"})

		case "progress_update":
			handleProgressUpdate(conn, hub, db, msg)

		default:
			send(conn, Message{Type: "error", Text: "unknown message type"})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("[TCP] connection error:", err)
	}
}

func handleProgressUpdate(conn net.Conn, hub *Hub, db *sql.DB, msg Message) {
	userID, username, err := auth.ValidateJWT(msg.Token)
	if err != nil {
		log.Println("[TCP] auth error:", err)
		send(conn, Message{Type: "error", Text: "invalid or missing token"})
		return
	}

	if msg.MangaID == "" {
		send(conn, Message{Type: "error", Text: "manga_id is required"})
		return
	}

	if msg.CurrentChapter < 0 {
		send(conn, Message{Type: "error", Text: "current_chapter cannot be negative"})
		return
	}

	if msg.Status == "" {
		msg.Status = "reading"
	}

	err = saveProgress(db, userID, msg.MangaID, msg.CurrentChapter, msg.Status)
	if err != nil {
		log.Println("[TCP] database error:", err)
		send(conn, Message{Type: "error", Text: "failed to save progress"})
		return
	}

	log.Printf(
		"[TCP] progress saved: user=%s username=%s manga=%s chapter=%d status=%s\n",
		userID,
		username,
		msg.MangaID,
		msg.CurrentChapter,
		msg.Status,
	)

	hub.Broadcast(Message{
		Type:           "progress_broadcast",
		UserID:         userID,
		MangaID:        msg.MangaID,
		CurrentChapter: msg.CurrentChapter,
		Status:         msg.Status,
		Text:           "progress updated",
	})
}

func saveProgress(db *sql.DB, userID, mangaID string, currentChapter int, status string) error {
	if db == nil {
		return errors.New("database is nil")
	}

	query := `
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id)
		DO UPDATE SET
			current_chapter = excluded.current_chapter,
			status = excluded.status,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := db.Exec(query, userID, mangaID, currentChapter, status)
	return err
}

func send(conn net.Conn, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("[TCP] send marshal error:", err)
		return
	}

	data = append(data, '\n')

	if _, err := conn.Write(data); err != nil {
		log.Println("[TCP] send error:", err)
	}
}
