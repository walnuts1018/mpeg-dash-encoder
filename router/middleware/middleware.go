package middleware

import (
	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/usecase"
)

type Middleware struct {
	adminToken string
	usecase    *usecase.Usecase
}

func NewMiddleware(adminToken config.AdminToken, usecase *usecase.Usecase) *Middleware {
	return &Middleware{
		adminToken: string(adminToken),
		usecase:    usecase,
	}
}
