package database

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mangahub/pkg/models"
)

type MangaDexResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Title       map[string]string `json:"title"`
			Description map[string]string `json:"description"`
			Status      string            `json:"status"`
			Tags        []struct {
				Attributes struct {
					Name map[string]string `json:"name"`
				} `json:"attributes"`
			} `json:"tags"`
		} `json:"attributes"`
	} `json:"data"`
}

func FetchMangaDex(limit int) ([]models.Manga, error) {
	if limit <= 0 {
		limit = 10
	}

	url := "https://api.mangadex.org/manga?limit=" + strconv.Itoa(limit)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var result MangaDexResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var mangaList []models.Manga

	for _, item := range result.Data {
		title := item.Attributes.Title["en"]
		if title == "" {
			for _, t := range item.Attributes.Title {
				title = t
				break
			}
		}

		description := item.Attributes.Description["en"]
		if description == "" {
			description = "No description"
		}

		var genres []string
		for _, tag := range item.Attributes.Tags {
			name := tag.Attributes.Name["en"]
			if name != "" {
				genres = append(genres, name)
			}
		}

		if len(genres) == 0 {
			genres = []string{"Unknown"}
		}

		description = strings.TrimSpace(description)
		if len(description) > 300 {
			description = description[:300] + "..."
		}

		// skip invalid data
		if item.ID == "" || title == "" {
			continue
		}

		mangaList = append(mangaList, models.Manga{
			ID:            item.ID,
			Title:         title,
			Author:        "Unknown",
			Genres:        genres,
			Status:        item.Attributes.Status,
			TotalChapters: 0,
			Description:   description,
		})
	}

	return mangaList, nil
}
