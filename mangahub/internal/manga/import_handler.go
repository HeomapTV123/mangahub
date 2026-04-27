package manga

import (
	"net/http"
	"strconv"

	"mangahub/pkg/database"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ImportFromMangaDex(c *gin.Context) {
	limitText := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitText)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be a positive number",
		})
		return
	}

	if limit > 50 {
		limit = 50
	}

	apiManga, err := database.FetchMangaDex(limit)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.SeedManga(h.Service.Repo.DB, apiManga)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save MangaDex data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "MangaDex import completed",
		"requested_limit": limit,
		"fetched_count":   len(apiManga),
	})
}
