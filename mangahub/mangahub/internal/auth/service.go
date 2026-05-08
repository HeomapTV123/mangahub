package auth

import (
	"database/sql"
	"errors"

	"mangahub/pkg/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Register(req models.RegisterRequest) error {
	if req.Username == "" || req.Password == "" {
		return errors.New("username and password are required")
	}

	var existingID string
	checkQuery := `SELECT id FROM users WHERE username = ?`
	err := s.DB.QueryRow(checkQuery, req.Username).Scan(&existingID)
	if err == nil {
		return errors.New("username already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userID := uuid.NewString()

	insertQuery := `
		INSERT INTO users (id, username, password_hash)
		VALUES (?, ?, ?)
	`
	_, err = s.DB.Exec(insertQuery, userID, req.Username, string(hashedPassword))
	return err
}

func (s *Service) Login(req models.LoginRequest) (string, error) {
	if req.Username == "" || req.Password == "" {
		return "", errors.New("username and password are required")
	}

	var userID string
	var passwordHash string

	query := `
		SELECT id, password_hash
		FROM users
		WHERE username = ?
	`
	err := s.DB.QueryRow(query, req.Username).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := GenerateJWT(userID, req.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
