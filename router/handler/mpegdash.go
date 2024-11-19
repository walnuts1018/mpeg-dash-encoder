package handler

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMediaFile(c *gin.Context) {
	mediaID := c.Param("media_id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "media_id is required"})
		return
	}

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	authorizedMediaIDs, err := h.getAuthorizedMediaIDs(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !slices.Contains(authorizedMediaIDs, mediaID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to access this media"})
		return
	}

	file, err := h.usecase.GetMediaFile(c.Request.Context(), mediaID, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get media file"})
		return
	}
	defer file.Close()

	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", file, map[string]string{
		"Cache-Control": "max-age=31536000",
	})
}
