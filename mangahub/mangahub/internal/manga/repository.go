package manga

import (
	"database/sql"
	"encoding/json"

	"mangahub/pkg/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAll() ([]models.Manga, error) {
	query := `
		SELECT id, title, author, genres, status, total_chapters, description
		FROM manga
		ORDER BY title ASC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangaList []models.Manga

	for rows.Next() {
		var m models.Manga
		var genresText string

		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Author,
			&genresText,
			&m.Status,
			&m.TotalChapters,
			&m.Description,
		)
		if err != nil {
			return nil, err
		}

		if genresText != "" {
			if err := json.Unmarshal([]byte(genresText), &m.Genres); err != nil {
				m.Genres = []string{}
			}
		}

		mangaList = append(mangaList, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mangaList, nil
}

func (r *Repository) GetByID(id string) (*models.Manga, error) {
	query := `
		SELECT id, title, author, genres, status, total_chapters, description
		FROM manga
		WHERE id = ?
	`

	var m models.Manga
	var genresText string

	err := r.DB.QueryRow(query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Author,
		&genresText,
		&m.Status,
		&m.TotalChapters,
		&m.Description,
	)
	if err != nil {
		return nil, err
	}

	if genresText != "" {
		if err := json.Unmarshal([]byte(genresText), &m.Genres); err != nil {
			m.Genres = []string{}
		}
	}

	return &m, nil
}

func (r *Repository) Create(m models.Manga) error {
	genresJSON, err := json.Marshal(m.Genres)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO manga (id, title, author, genres, status, total_chapters, description)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.DB.Exec(
		query,
		m.ID,
		m.Title,
		m.Author,
		string(genresJSON),
		m.Status,
		m.TotalChapters,
		m.Description,
	)

	return err
}

func (r *Repository) Update(id string, m models.Manga) error {
	genresJSON, err := json.Marshal(m.Genres)
	if err != nil {
		return err
	}

	query := `
	UPDATE manga
	SET title = ?, author = ?, genres = ?, status = ?, total_chapters = ?, description = ?
	WHERE id = ?
	`

	result, err := r.DB.Exec(
		query,
		m.Title,
		m.Author,
		string(genresJSON),
		m.Status,
		m.TotalChapters,
		m.Description,
		id,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	query := `DELETE FROM manga WHERE id = ?`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
