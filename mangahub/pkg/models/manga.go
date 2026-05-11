package models

type Manga struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status"`
	TotalChapters int      `json:"total_chapters"`
	Description   string   `json:"description"`
	CoverURL      string   `json:"cover_url,omitempty"`
}

type CreateMangaRequest struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status"`
	TotalChapters int      `json:"total_chapters"`
	Description   string   `json:"description"`
}

type SearchFilters struct {
	Keyword     string `json:"keyword"`
	Genre       string `json:"genre"`
	Status      string `json:"status"`
	MinChapters int    `json:"min_chapters"`
	MaxChapters int    `json:"max_chapters"`
	SortBy      string `json:"sort_by"`
}
