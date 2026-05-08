package database

import (
	"encoding/json"
	"os"

	"mangahub/pkg/models"
)

func LoadMangaJSON(filePath string) ([]models.Manga, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var mangaList []models.Manga
	if err := json.Unmarshal(data, &mangaList); err != nil {
		return nil, err
	}

	return mangaList, nil
}
