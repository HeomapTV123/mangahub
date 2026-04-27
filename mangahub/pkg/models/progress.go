package models

type UserProgress struct {
	UserID         string `json:"user_id"`
	MangaID        string `json:"manga_id"`
	CurrentChapter int    `json:"current_chapter"`
	Status         string `json:"status"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type AddLibraryRequest struct {
	MangaID string `json:"manga_id"`
	Status  string `json:"status"`
}

type UpdateProgressRequest struct {
	MangaID        string `json:"manga_id"`
	CurrentChapter int    `json:"current_chapter"`
	Status         string `json:"status"`
}
