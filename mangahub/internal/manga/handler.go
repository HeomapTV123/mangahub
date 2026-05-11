package manga

import (
	"mangahub/internal/udp"
	"mangahub/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service   *Service
	udpServer *udp.NotificationServer
}

func NewHandler(service *Service, udpServer *udp.NotificationServer) *Handler {
	return &Handler{
		Service:   service,
		udpServer: udpServer,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	mangaList, err := h.Service.GetAllManga()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch manga list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"manga": mangaList,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	manga, err := h.Service.GetMangaByID(id)
	if err != nil {
		if err.Error() == "manga not found" || err.Error() == "manga id is required" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch manga",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"manga": manga,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req models.CreateMangaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.Service.CreateManga(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	notif := udp.Notification{
		Type:      "new_chapter",
		MangaID:   req.ID,
		Message:   req.Title + " updated!",
		Timestamp: time.Now().Unix(),
	}

	udp.Broadcast(h.udpServer, notif)

	c.JSON(http.StatusCreated, gin.H{
		"message": "manga created",
	})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.CreateMangaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.Service.UpdateManga(id, req)
	if err != nil {
		if err.Error() == "manga not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "manga updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.DeleteManga(id)
	if err != nil {
		if err.Error() == "manga not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "manga deleted"})
}

func (h *Handler) Search(c *gin.Context) {
	minChapters := 0
	maxChapters := 0
	if minText := strings.TrimSpace(c.Query("min_chapters")); minText != "" {
		value, err := strconv.Atoi(minText)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "min_chapters must be a number",
			})
			return
		}
		minChapters = value
	}

	if maxText := strings.TrimSpace(c.Query("max_chapters")); maxText != "" {
		value, err := strconv.Atoi(maxText)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "max_chapters must be a number",
			})
			return
		}
		maxChapters = value
	}

	filters := models.SearchFilters{
		Keyword:     c.Query("keyword"),
		Genre:       c.Query("genre"),
		Status:      c.Query("status"),
		MinChapters: minChapters,
		MaxChapters: maxChapters,
		SortBy:      c.DefaultQuery("sort_by", "title"),
	}

	results, err := h.Service.SearchManga(filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"manga":   results,
	})
}
