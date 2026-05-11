package manga

import (
	"database/sql"
	"errors"
	"strings"

	"mangahub/pkg/models"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllManga() ([]models.Manga, error) {
	return s.Repo.GetAll()
}

func (s *Service) GetMangaByID(id string) (*models.Manga, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return nil, errors.New("manga id is required")
	}

	manga, err := s.Repo.GetByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.New("manga not found")
	}

	return manga, err
}

func (s *Service) CreateManga(req models.CreateMangaRequest) error {
	req.ID = strings.TrimSpace(req.ID)
	req.Title = strings.TrimSpace(req.Title)
	req.Author = strings.TrimSpace(req.Author)
	req.Status = strings.TrimSpace(req.Status)
	req.Description = strings.TrimSpace(req.Description)

	if req.ID == "" || req.Title == "" {
		return errors.New("id and title are required")
	}

	if req.TotalChapters < 0 {
		return errors.New("total_chapters cannot be negative")
	}

	if req.Status == "" {
		req.Status = "ongoing"
	}

	if !isValidMangaStatus(req.Status) {
		return errors.New("status must be one of: ongoing, completed, hiatus, cancelled")
	}

	for i, genre := range req.Genres {
		req.Genres[i] = strings.TrimSpace(genre)
	}

	m := models.Manga{
		ID:            req.ID,
		Title:         req.Title,
		Author:        req.Author,
		Genres:        req.Genres,
		Status:        req.Status,
		TotalChapters: req.TotalChapters,
		Description:   req.Description,
	}

	return s.Repo.Create(m)
}

func (s *Service) UpdateManga(id string, req models.CreateMangaRequest) error {
	id = strings.TrimSpace(id)
	req.Title = strings.TrimSpace(req.Title)
	req.Author = strings.TrimSpace(req.Author)
	req.Status = strings.TrimSpace(req.Status)
	req.Description = strings.TrimSpace(req.Description)

	if id == "" {
		return errors.New("id is required")
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.TotalChapters < 0 {
		return errors.New("total_chapters cannot be negative")
	}

	if req.Status == "" {
		req.Status = "ongoing"
	}

	if !isValidMangaStatus(req.Status) {
		return errors.New("status must be one of: ongoing, completed, hiatus, cancelled")
	}

	for i, genre := range req.Genres {
		req.Genres[i] = strings.TrimSpace(genre)
	}

	m := models.Manga{
		ID:            id,
		Title:         req.Title,
		Author:        req.Author,
		Genres:        req.Genres,
		Status:        req.Status,
		TotalChapters: req.TotalChapters,
		Description:   req.Description,
	}

	err := s.Repo.Update(id, m)
	if err == sql.ErrNoRows {
		return errors.New("manga not found")
	}

	return err
}

func (s *Service) DeleteManga(id string) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return errors.New("id is required")
	}

	err := s.Repo.Delete(id)
	if err == sql.ErrNoRows {
		return errors.New("manga not found")
	}

	return err
}

func isValidMangaStatus(status string) bool {
	switch status {
	case "ongoing", "completed", "hiatus", "cancelled":
		return true
	default:
		return false
	}
}

func (s *Service) SearchManga(filters models.SearchFilters) ([]models.Manga, error) {
	filters.Keyword = strings.TrimSpace(filters.Keyword)
	filters.Genre = strings.TrimSpace(filters.Genre)
	filters.Status = strings.TrimSpace(filters.Status)
	filters.SortBy = strings.TrimSpace(filters.SortBy)

	if filters.MinChapters < 0 {
		return nil, errors.New("min_chapters cannot be negative")
	}

	if filters.MaxChapters < 0 {
		return nil, errors.New("max_chapters cannot be negative")
	}

	if filters.MinChapters > 0 && filters.MaxChapters > 0 && filters.MinChapters > filters.MaxChapters {
		return nil, errors.New("min_chapters cannot be greater than max_chapters")
	}

	if filters.Status != "" && !isValidMangaStatus(filters.Status) {
		return nil, errors.New("status must be one of: ongoing, completed, hiatus, cancelled")
	}

	if filters.SortBy != "" && !isValidSortBy(filters.SortBy) {
		return nil, errors.New("sort_by must be one of: title, title_desc, chapters, chapters_desc, recent")
	}

	return s.Repo.Search(filters)
}

func isValidSortBy(sortBy string) bool {
	switch sortBy {
	case "title", "title_desc", "chapters", "chapters_desc", "recent":
		return true
	default:
		return false
	}
}
