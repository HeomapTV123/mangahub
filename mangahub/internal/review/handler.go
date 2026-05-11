package review

import (
	"net/http"

	"mangahub/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) UpsertReview(c *gin.Context) {
	userID := c.GetString("user_id")
	mangaID := c.Param("id")

	var req models.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.Service.UpsertReview(userID, mangaID, req)
	if err != nil {
		status := http.StatusBadRequest

		if err.Error() == "manga not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "review saved",
	})
}

func (h *Handler) GetReviewsByMangaID(c *gin.Context) {
	mangaID := c.Param("id")

	reviews, err := h.Service.GetReviewsByMangaID(mangaID)
	if err != nil {
		status := http.StatusBadRequest

		if err.Error() == "manga not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(reviews),
		"reviews": reviews,
	})
}

func (h *Handler) GetReviewSummary(c *gin.Context) {
	mangaID := c.Param("id")

	summary, err := h.Service.GetReviewSummary(mangaID)
	if err != nil {
		status := http.StatusBadRequest

		if err.Error() == "manga not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}

func (h *Handler) DeleteReview(c *gin.Context) {
	userID := c.GetString("user_id")
	mangaID := c.Param("id")

	err := h.Service.DeleteReview(userID, mangaID)
	if err != nil {
		status := http.StatusBadRequest

		if err.Error() == "review not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "review deleted",
	})
}
