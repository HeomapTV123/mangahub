package review

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"mangahub/pkg/models"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) UpsertReview(userID string, mangaID string, req models.ReviewRequest) error {
	userID = strings.TrimSpace(userID)
	mangaID = strings.TrimSpace(mangaID)
	req.Text = strings.TrimSpace(req.Text)

	if userID == "" {
		return errors.New("user id is required")
	}

	if mangaID == "" {
		return errors.New("manga id is required")
	}

	exists, err := s.Repo.MangaExists(mangaID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("manga not found")
	}

	if req.Rating < 1 || req.Rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}

	if len(req.Text) > 1000 {
		return errors.New("review text cannot be longer than 1000 characters")
	}

	review := models.Review{
		UserID:    userID,
		MangaID:   mangaID,
		Rating:    req.Rating,
		Text:      req.Text,
		Timestamp: time.Now().Unix(),
		Helpful:   0,
	}

	return s.Repo.UpsertReview(review)
}

func (s *Service) GetReviewsByMangaID(mangaID string) ([]models.Review, error) {
	mangaID = strings.TrimSpace(mangaID)

	if mangaID == "" {
		return nil, errors.New("manga id is required")
	}

	exists, err := s.Repo.MangaExists(mangaID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("manga not found")
	}

	return s.Repo.GetReviewsByMangaID(mangaID)
}

func (s *Service) GetReviewSummary(mangaID string) (*models.ReviewSummary, error) {
	mangaID = strings.TrimSpace(mangaID)

	if mangaID == "" {
		return nil, errors.New("manga id is required")
	}

	exists, err := s.Repo.MangaExists(mangaID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("manga not found")
	}

	return s.Repo.GetReviewSummary(mangaID)
}

func (s *Service) DeleteReview(userID string, mangaID string) error {
	userID = strings.TrimSpace(userID)
	mangaID = strings.TrimSpace(mangaID)

	if userID == "" {
		return errors.New("user id is required")
	}

	if mangaID == "" {
		return errors.New("manga id is required")
	}

	err := s.Repo.DeleteReview(userID, mangaID)
	if err == sql.ErrNoRows {
		return errors.New("review not found")
	}

	return err
}
