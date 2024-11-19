package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUserToken(c *gin.Context) {
	var req struct {
		MediaIDs []string `json:"media_ids"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	token, err := h.usecase.CreateUserToken(req.MediaIDs)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create token",
		})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
