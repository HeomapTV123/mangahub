package manga

import (
	"database/sql"
	"errors"

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
	if id == "" {
		return nil, errors.New("manga id is required")
	}

	manga, err := s.Repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("manga not found")
		}
		return nil, err
	}

	return manga, nil
}

func (s *Service) CreateManga(req models.CreateMangaRequest) error {
	if req.ID == "" || req.Title == "" {
		return errors.New("id and title are required")
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
	if id == "" {
		return errors.New("id is required")
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
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("manga not found")
		}
		return err
	}

	return nil
}

func (s *Service) DeleteManga(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	err := s.Repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("manga not found")
		}
		return err
	}

	return nil
}
