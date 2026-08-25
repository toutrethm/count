package handler

import (
	"count/backend/internal/middleware"
	"count/backend/pkg/auth"

	"github.com/gin-gonic/gin"
)

func currentClaims(c *gin.Context) (*auth.Claims, bool) {
	return middleware.ClaimsFromContext(c)
}
