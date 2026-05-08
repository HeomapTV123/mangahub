package websocket

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleConnection(hub *ChatHub, conn *websocket.Conn, username string) {
	hub.Register <- ClientConnection{
		Conn:     conn,
		Username: username,
	}

	defer func() {
		hub.Unregister <- conn
	}()

	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("[WS] Read error:", err)
			break
		}

		hub.Broadcast <- msg
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *ChatHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		go HandleConnection(hub, conn, username)
	}
}
