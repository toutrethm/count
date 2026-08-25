package handler

import (
	"net/http"
	"strings"

	"count/backend/internal/model"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Phone = strings.TrimSpace(req.Phone)
	req.StationRole = strings.TrimSpace(req.StationRole)
	if req.Username == "" || req.Phone == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing fields"})
		return
	}

	var count int64
	if err := h.DB.Model(&model.User{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "phone already exists"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user := model.User{
		Username:     req.Username,
		Phone:        req.Phone,
		PasswordHash: string(passwordHash),
		Role:         "worker",
		StationRole:  req.StationRole,
		Status:       1,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user model.User
	if err := h.DB.Where("phone = ? OR username = ?", req.Account, req.Account).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid account or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid account or password"})
		return
	}

	token, err := tokenForUser(h.JWTSecret, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *Handler) Me(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	var user model.User
	if err := h.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	processes, _ := h.assignedProcesses(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"user":      user,
		"processes": processes,
	})
}
