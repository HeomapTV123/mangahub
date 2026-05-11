package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"mangahub/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *ChatHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(c.Query("token"))
		mangaID := strings.TrimSpace(c.Query("manga_id"))

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token is required",
			})
			return
		}

		if mangaID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "manga_id is required",
			})
			return
		}

		userID, username, err := auth.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("[WS] Upgrade error:", err)
			return
		}

		go HandleConnection(hub, conn, userID, username, mangaID)
	}
}

func HandleConnection(hub *ChatHub, conn *websocket.Conn, userID string, username string, mangaID string) {
	client := ClientConnection{
		Conn:     conn,
		UserID:   userID,
		Username: username,
		MangaID:  mangaID,
	}

	hub.Register <- client

	defer func() {
		hub.Unregister <- client
	}()

	for {
		var msg ChatMessage

		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("[WS] Read error:", err)
			break
		}

		msg.UserID = userID
		msg.Username = username
		msg.MangaID = mangaID

		if msg.Timestamp == 0 {
			msg.Timestamp = time.Now().Unix()
		}

		if strings.TrimSpace(msg.Message) == "" {
			continue
		}

		hub.Broadcast <- RoomMessage{
			MangaID: mangaID,
			Message: msg,
		}
	}
}
