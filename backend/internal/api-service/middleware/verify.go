package middleware

import (
	"backend/internal/errors"
	"backend/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrNeedVerified errors.SimpleError = ""
)

func (s *Service) VerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*schema.User)

		if !user.VerifiedAt.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrNeedVerified)
			return
		}

		c.Next()
	}
}
