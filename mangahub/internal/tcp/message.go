package tcp

type Message struct {
	Type           string `json:"type"`
	Token          string `json:"token,omitempty"`
	MangaID        string `json:"manga_id,omitempty"`
	CurrentChapter int    `json:"current_chapter,omitempty"`
	Status         string `json:"status,omitempty"`
	Text           string `json:"text,omitempty"`
	UserID         string `json:"user_id,omitempty"`
}
