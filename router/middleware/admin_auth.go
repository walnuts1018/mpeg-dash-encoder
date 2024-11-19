package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if got, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer "); ok && got == m.adminToken {
			c.Next()
			return
		} else {
			c.JSON(401, gin.H{
				"error": "authorization failed",
			})
			c.Abort()
			return
		}
	}
}
