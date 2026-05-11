package review

import (
	"database/sql"

	"mangahub/pkg/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) UpsertReview(review models.Review) error {
	query := `
		INSERT INTO reviews (user_id, manga_id, rating, text, timestamp, helpful)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, manga_id)
		DO UPDATE SET
			rating = excluded.rating,
			text = excluded.text,
			timestamp = excluded.timestamp
	`

	_, err := r.DB.Exec(
		query,
		review.UserID,
		review.MangaID,
		review.Rating,
		review.Text,
		review.Timestamp,
		review.Helpful,
	)

	return err
}

func (r *Repository) GetReviewsByMangaID(mangaID string) ([]models.Review, error) {
	query := `
		SELECT r.user_id, u.username, r.manga_id, r.rating, r.text, r.timestamp, r.helpful
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.manga_id = ?
		ORDER BY r.timestamp DESC
	`

	rows, err := r.DB.Query(query, mangaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review

	for rows.Next() {
		var review models.Review

		err := rows.Scan(
			&review.UserID,
			&review.Username,
			&review.MangaID,
			&review.Rating,
			&review.Text,
			&review.Timestamp,
			&review.Helpful,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}

func (r *Repository) GetReviewSummary(mangaID string) (*models.ReviewSummary, error) {
	query := `
		SELECT COUNT(*), COALESCE(AVG(rating), 0)
		FROM reviews
		WHERE manga_id = ?
	`

	var count int
	var average float64

	err := r.DB.QueryRow(query, mangaID).Scan(&count, &average)
	if err != nil {
		return nil, err
	}

	return &models.ReviewSummary{
		MangaID:       mangaID,
		AverageRating: average,
		ReviewCount:   count,
	}, nil
}

func (r *Repository) DeleteReview(userID string, mangaID string) error {
	query := `
		DELETE FROM reviews
		WHERE user_id = ? AND manga_id = ?
	`

	result, err := r.DB.Exec(query, userID, mangaID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) MangaExists(mangaID string) (bool, error) {
	query := `SELECT id FROM manga WHERE id = ?`

	var id string
	err := r.DB.QueryRow(query, mangaID).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
