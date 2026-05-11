package models

type Review struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username,omitempty"`
	MangaID   string `json:"manga_id"`
	Rating    int    `json:"rating"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
	Helpful   int    `json:"helpful"`
}

type ReviewRequest struct {
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

type ReviewSummary struct {
	MangaID       string  `json:"manga_id"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int     `json:"review_count"`
}
