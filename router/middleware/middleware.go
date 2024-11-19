package middleware

import (
	"github.com/walnuts1018/mpeg_dash-encoder/config"
	"github.com/walnuts1018/mpeg_dash-encoder/usecase"
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
