package udp

import "net"

type NotificationServer struct {
	Port    string
	Clients []net.UDPAddr
}
type Notification struct {
	Type      string `json:"type"`
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
