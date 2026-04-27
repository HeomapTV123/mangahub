package user

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

func (r *Repository) AddToLibrary(userID string, req models.AddLibraryRequest) error {
	query := `
	INSERT OR REPLACE INTO user_progress 
	(user_id, manga_id, current_chapter, status, updated_at)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	_, err := r.DB.Exec(query, userID, req.MangaID, 0, req.Status)
	return err
}

func (r *Repository) GetLibrary(userID string) ([]models.UserProgress, error) {
	query := `
	SELECT user_id, manga_id, current_chapter, status, updated_at
	FROM user_progress
	WHERE user_id = ?
	ORDER BY updated_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.UserProgress

	for rows.Next() {
		var p models.UserProgress

		err := rows.Scan(
			&p.UserID,
			&p.MangaID,
			&p.CurrentChapter,
			&p.Status,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, p)
	}

	return list, rows.Err()
}

func (r *Repository) UpdateProgress(userID string, req models.UpdateProgressRequest) error {
	query := `
	UPDATE user_progress
	SET current_chapter = ?, status = ?, updated_at = CURRENT_TIMESTAMP
	WHERE user_id = ? AND manga_id = ?
	`

	result, err := r.DB.Exec(query, req.CurrentChapter, req.Status, userID, req.MangaID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		insertQuery := `
		INSERT INTO user_progress 
		(user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		`

		_, err = r.DB.Exec(insertQuery, userID, req.MangaID, req.CurrentChapter, req.Status)
		return err
	}

	return nil
}
