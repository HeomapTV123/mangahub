package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"mangahub/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

func InitSQLite(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	fmt.Println("SQLite connected successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	mangaTable := `
	CREATE TABLE IF NOT EXISTS manga (
		id TEXT PRIMARY KEY,
		title TEXT,
		author TEXT,
		genres TEXT,
		status TEXT,
		total_chapters INTEGER,
		description TEXT
	);`

	userProgressTable := `
	CREATE TABLE IF NOT EXISTS user_progress (
		user_id TEXT,
		manga_id TEXT,
		current_chapter INTEGER,
		status TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id)
	);`

	reviewTable := `
	CREATE TABLE IF NOT EXISTS reviews (
	user_id TEXT NOT NULL,
	manga_id TEXT NOT NULL,
	rating INTEGER NOT NULL,
	text TEXT,
	timestamp INTEGER NOT NULL,
	helpful INTEGER DEFAULT 0,
	PRIMARY KEY (user_id, manga_id)
	);`

	if _, err := db.Exec(reviewTable); err != nil {
		return err
	}

	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	if _, err := db.Exec(mangaTable); err != nil {
		return err
	}

	if _, err := db.Exec(userProgressTable); err != nil {
		return err
	}

	return nil
}

func SeedManga(db *sql.DB, mangaList []models.Manga) error {
	for _, manga := range mangaList {
		genresJSON, err := json.Marshal(manga.Genres)
		if err != nil {
			return err
		}

		query := `
		INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`

		_, err = db.Exec(
			query,
			manga.ID,
			manga.Title,
			manga.Author,
			string(genresJSON),
			manga.Status,
			manga.TotalChapters,
			manga.Description,
		)
		if err != nil {
			return err
		}
	}

	fmt.Println("Manga seed completed")
	return nil
}

func IsMangaTableEmpty(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM manga").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
