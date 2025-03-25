package handler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/usecase"
)

type Handler struct {
	config  config.Config
	usecase *usecase.Usecase
}

func NewHandler(config config.Config, usecase *usecase.Usecase) (Handler, error) {
	return Handler{
		config,
		usecase,
	}, nil
}

func (h *Handler) getAuthorizedMediaIDs(c *gin.Context) ([]string, error) {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok {
		return nil, errors.New("authorization header is missing")
	}
	return h.usecase.GetMediaIDsFromToken(token)
}
