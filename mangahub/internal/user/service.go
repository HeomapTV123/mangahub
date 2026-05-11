package user

import (
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

func (s *Service) AddToLibrary(userID string, req models.AddLibraryRequest) error {
	userID = strings.TrimSpace(userID)
	req.MangaID = strings.TrimSpace(req.MangaID)
	req.Status = strings.TrimSpace(req.Status)

	if userID == "" {
		return errors.New("user id is required")
	}

	if req.MangaID == "" {
		return errors.New("manga id is required")
	}

	if req.Status == "" {
		req.Status = "plan_to_read"
	}

	if !isValidProgressStatus(req.Status) {
		return errors.New("status must be one of: plan_to_read, reading, completed, dropped")
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
	userID = strings.TrimSpace(userID)
	req.MangaID = strings.TrimSpace(req.MangaID)
	req.Status = strings.TrimSpace(req.Status)

	if userID == "" {
		return errors.New("user id is required")
	}

	if req.MangaID == "" {
		return errors.New("manga id is required")
	}

	if req.CurrentChapter < 0 {
		return errors.New("current_chapter cannot be negative")
	}

	if req.Status == "" {
		req.Status = "reading"
	}

	if !isValidProgressStatus(req.Status) {
		return errors.New("status must be one of: plan_to_read, reading, completed, dropped")
	}

	return s.Repo.UpdateProgress(userID, req)
}

func isValidProgressStatus(status string) bool {
	switch status {
	case "plan_to_read", "reading", "completed", "dropped":
		return true
	default:
		return false
	}
}
