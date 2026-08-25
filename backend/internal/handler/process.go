package handler

import (
	"net/http"
	"strings"

	"count/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListProcesses(c *gin.Context) {
	processes, err := h.allProcesses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": processes})
}

func (h *Handler) ListMyProcesses(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	processes, err := h.assignedProcesses(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": processes})
}

func (h *Handler) CreateProcess(c *gin.Context) {
	var req CreateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	process := model.Process{
		Code:   strings.TrimSpace(req.Code),
		Name:   strings.TrimSpace(req.Name),
		Sort:   req.Sort,
		Status: 1,
	}
	if err := h.DB.Create(&process).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": process})
}
