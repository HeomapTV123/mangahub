package user

import (
	"errors"

	"mangahub/pkg/models"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) AddToLibrary(userID string, req models.AddLibraryRequest) error {
	if userID == "" {
		return errors.New("user id is required")
	}
	if req.MangaID == "" {
		return errors.New("manga id is required")
	}
	if req.Status == "" {
		req.Status = "plan_to_read"
	}

	return s.Repo.AddToLibrary(userID, req)
}

func (s *Service) GetLibrary(userID string) ([]models.UserProgress, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	return s.Repo.GetLibrary(userID)
}

func (s *Service) UpdateProgress(userID string, req models.UpdateProgressRequest) error {
	if userID == "" {
		return errors.New("user id is required")
	}
	if req.MangaID == "" {
		return errors.New("manga id is required")
	}
	if req.Status == "" {
		req.Status = "reading"
	}

	return s.Repo.UpdateProgress(userID, req)
}
